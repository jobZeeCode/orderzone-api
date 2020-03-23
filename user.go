package main

import (
	"context"
	"encoding/json"
	"errors"
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

const db = "customers"

//UserHandler : manage database
func UserHandler(w http.ResponseWriter, r *http.Request) {
	client, ctx = Conect()
	w.Header().Set("content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, content-type")
	defer client.Close()

	switch r.Method {
	case "GET":
		//ReadData
		id := r.URL.Query()["id"]
		if len(id) != 0 {
			//have query
			data, err := serchDataFromID(id[0], fetchData())
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"msg": "not found"})
			} else {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(data)
			}
		} else {
			//Not have query
			data := fetchData()
			if len(data) > 0 {
				//
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(fetchData())
			} else {
				json.NewEncoder(w).Encode(map[string]interface{}{"msg": "not found"})
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

func serchDataFromID(id string, data []map[string]interface{}) (map[string]interface{}, error) {
	for i, v := range data {
		if id == v["ID"] {
			return data[i], nil
		}
	}
	return nil, errors.New("not Found data")
}
