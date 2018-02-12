package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/mostafa-asg/finch/core"
	"github.com/mostafa-asg/finch/generator/base62"
	"github.com/mostafa-asg/finch/http/model"
	"github.com/mostafa-asg/finch/service/registrator"
	"github.com/mostafa-asg/finch/storage/cassandra"
	"github.com/mostafa-asg/finch/storage/sqlite"
	"github.com/prometheus/client_golang/prometheus"
	config "github.com/spf13/viper"
	"github.com/teris-io/shortid"
)

var storage core.Storage
var generator core.Generator

func main() {

	var configPath string
	flag.StringVar(&configPath, "config", "configs/finch.yml", "config file path")
	flag.Parse()

	if !isFileExists(configPath) {
		log.Fatal("Config file does not exist", configPath)
	}

	configDir, configFilename := filepath.Split(configPath)
	configFilename = configFilename[0:strings.LastIndexByte(configFilename, byte('.'))]

	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.SetConfigName(configFilename)
	config.AddConfigPath(configDir)
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading config file", err)
	}

	storage = instantiateStorage()
	generator = base62.NewConcurrent()

	router := mux.NewRouter()
	router.HandleFunc("/get/{id}", getHandler()).Methods("GET")
	router.HandleFunc("/hash", hashHandler()).Methods("POST")
	router.Handle("/metrics", prometheus.Handler())

	serviceID, err := shortid.Generate()
	if err != nil {
		log.Fatal("Error in creating short uuid ", err)
	}

	serverAddress := config.GetString("server.address")
	port := config.GetInt("server.port")

	server := &http.Server{
		Addr:    serverAddress + ":" + strconv.Itoa(port),
		Handler: router,
	}
	go func() {

		hostIP := resolveHostIp()
		if hostIP == "" {
			log.Fatal("could not find host IP address")
		}

		go func() {
			for {
				err = registrator.NewConsulServiceDiscovery().Register("finch-REST", serviceID, hostIP, port)
				if err != nil {
					log.Println("Unable to register service ", err)
					//Try one second later
					time.Sleep(1 * time.Second)
				} else {
					break
				}
			}
		}()

		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}

	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	deregisterService(serviceID)
	server.Shutdown(context.Background())
}

func resolveHostIp() string {

	netInterfaceAddresses, err := net.InterfaceAddrs()

	if err != nil {
		log.Println(err)
		return ""
	}

	for _, netInterfaceAddress := range netInterfaceAddresses {

		networkIp, ok := netInterfaceAddress.(*net.IPNet)

		if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {

			ip := networkIp.IP.String()

			log.Println("Resolved Host IP: " + ip)

			return ip
		}
	}
	return ""
}

func deregisterService(serviceID string) {
	err := registrator.NewConsulServiceDiscovery().Deregister(serviceID)
	if err != nil {
		log.Println("Unable to deregister service ", err)
	}
}

func instantiateStorage() core.Storage {

	storage := config.GetString("storage.type")
	switch storage {
	case "sqlite":
		return sqlite.New()
	case "cassandra":
		return cassandra.New()
	}

	log.Fatal(fmt.Sprintf("Unknown storage type : %s", storage))
	return nil
}

// Exists reports whether the named file or directory exists.
func isFileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func hashHandler() func(http.ResponseWriter, *http.Request) {

	count := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "api_http_hash_request_total",
		Help: "Total hash requests count",
	})

	collisions := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "api_http_collisions_total",
		Help: "Total database ID collisions count",
	})

	failer := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "api_http_hash_fail_total",
		Help: "Total number of unanswered hash requests",
	})

	prometheus.MustRegister(count)
	prometheus.MustRegister(collisions)
	prometheus.MustRegister(failer)

	return func(w http.ResponseWriter, r *http.Request) {

		count.Inc()

		var request model.HashRequest
		json.NewDecoder(r.Body).Decode(&request)

		var newID string
		try := 0
		for {
			newID = generator.GenerateID()
			err := storage.Put(newID, request.Url)
			if err == nil {
				break
			}

			//maybe we had generated duplicate id
			//so we try again
			try++
			collisions.Inc()
			if try >= 20 {
				//maybe almost all ID has been taken
				//or maybe the storage system has been down
				failer.Inc()
				json.NewEncoder(w).Encode(model.HashResponse{
					ErrorMessage: "Sorry , we cannot serve your request right now",
				})
				return
			}
		}

		json.NewEncoder(w).Encode(model.HashResponse{
			TinyUrl: newID,
		})

	}
}

func getHandler() func(http.ResponseWriter, *http.Request) {

	count := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_get_request_total",
		Help: "Total get requests count",
	}, []string{"type"})

	prometheus.MustRegister(count)

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id := params["id"]

		url, err := storage.Get(id)
		if err != nil {
			count.WithLabelValues("miss").Inc()
			json.NewEncoder(w).Encode(model.GetResponse{
				Found: false,
			})
			return
		}

		count.WithLabelValues("hit").Inc()
		json.NewEncoder(w).Encode(model.GetResponse{
			Found: true,
			Url:   url,
		})
	}
}
