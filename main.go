package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mostafa-asg/finch/core"
	"github.com/mostafa-asg/finch/generator/base62"
	"github.com/mostafa-asg/finch/storage/sqlite"
)

var storage core.Storage
var generator core.Generator

func main() {

	storage = sqlite.New()
	generator = base62.NewConcurrent()

	router := mux.NewRouter()
	router.HandleFunc("/get/{id}", getHandler).Methods("GET")
	router.HandleFunc("/hash", hashHandler).Methods("POST")
	http.ListenAndServe(":8585", router)
}

type hashResponse struct {
	TinyUrl string `json:"tiny"`
}
type hashRequest struct {
	Url string `json:"url"`
}

func hashHandler(w http.ResponseWriter, r *http.Request) {

	var request hashRequest
	json.NewDecoder(r.Body).Decode(&request)

	var newID string
	for {
		newID = generator.GenerateID()
		inserted, err := storage.Put(newID, request.Url)
		if err != nil && err != sqlite.ErrUnique {
			log.Fatal(err)
		}
		if inserted {
			break
		}
	}

	json.NewEncoder(w).Encode(hashResponse{
		TinyUrl: newID,
	})

}

type getResponse struct {
	Found bool   `json:"found"`
	Url   string `json:"url,omitempty"`
}

func getHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	url, err := storage.Get(id)
	if err != nil {
		json.NewEncoder(w).Encode(getResponse{
			Found: false,
		})
		return
	}

	json.NewEncoder(w).Encode(getResponse{
		Found: true,
		Url:   url,
	})
}
