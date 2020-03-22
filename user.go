package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type userStruct struct {
	ID       interface{}
	Name     interface{}
	Lastname interface{}
	Email    interface{}
	Password interface{}
	Tel      interface{}
}

var client *firestore.Client
var ctx context.Context

const db = "users"

//UserHandler : manage database
func UserHandler(w http.ResponseWriter, r *http.Request) {
	client, ctx = Conect()
	w.Header().Set("content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, content-type")
	defer client.Close()

	switch r.Method {
	case "GET":
		id := r.URL.Query()["id"]
		if len(id) > 1 {
			//have query
			json.NewEncoder(w).Encode(fetchData())
		} else {
			//Not have query
			data := fetchData()
			if len(data) > 0 {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(fetchData())
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		}
	}
}

func fetchData() []map[string]interface{} {
	var data []map[string]interface{}
	iter := client.Collection(db).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate : %v ", err)
		} else {
			data = append(data, doc.Data())
		}

	}
	return data
}
