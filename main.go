package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"

	"github.com/mostafa-asg/finch/core"
	"github.com/mostafa-asg/finch/generator/base62"
	"github.com/mostafa-asg/finch/storage/cassandra"
	"github.com/mostafa-asg/finch/storage/sqlite"
	"github.com/prometheus/client_golang/prometheus"
	config "github.com/spf13/viper"
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
	http.ListenAndServe(config.GetString("server.bind"), router)
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

type hashResponse struct {
	TinyUrl      string `json:"tiny,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}
type hashRequest struct {
	Url string `json:"url"`
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

	return func(w http.ResponseWriter, r *http.Request) {

		count.Inc()

		var request hashRequest
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
				json.NewEncoder(w).Encode(hashResponse{
					ErrorMessage: "Sorry , we cannot serve your request right now",
				})
				return
			}
		}

		json.NewEncoder(w).Encode(hashResponse{
			TinyUrl: newID,
		})

	}
}

type getResponse struct {
	Found bool   `json:"found"`
	Url   string `json:"url,omitempty"`
}

func getHandler() func(http.ResponseWriter, *http.Request) {

	count := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_get_request",
		Help: "Total get requests count",
	}, []string{"type"})

	prometheus.MustRegister(count)

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id := params["id"]

		url, err := storage.Get(id)
		if err != nil {
			count.WithLabelValues("miss").Inc()
			json.NewEncoder(w).Encode(getResponse{
				Found: false,
			})
			return
		}

		count.WithLabelValues("hit").Inc()
		json.NewEncoder(w).Encode(getResponse{
			Found: true,
			Url:   url,
		})
	}
}
