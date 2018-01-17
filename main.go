package main

import (
	"encoding/json"
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
	TinyUrl string `json:"tiny,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}
type hashRequest struct {
	Url string `json:"url"`
}

func hashHandler(w http.ResponseWriter, r *http.Request) {
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
		if try>=20 {
			//maybe almost ID has been taken
			//or maybe the storage system has been down
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
