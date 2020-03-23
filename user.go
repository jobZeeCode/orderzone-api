package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
		//read data
		id := r.URL.Query()["id"]
		if len(id) != 0 {
			//have query
			data, err := serchData(id[0], "ID", fetchData())
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
	case "POST":
		//insert data
		var user userStruct
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Fatalf("Fail Decoder : %v ", err)
		}
		ref := client.Collection(db).NewDoc()
		user.ID = ref.ID
		doc, err := ref.Set(ctx, user)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"msg": "add fail : ", "data": doc})
		} else {
			json.NewEncoder(w).Encode(user)
		}
	case "PUT":
		//Edit data
		var user userStruct
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Fatalf("Failt Decoder : %v ", err)
		}
		doc, err := client.Collection(db).Doc(fmt.Sprintf("%v", user.ID)).Set(ctx, user)
		if err != nil {
			log.Fatalf("Fail Decoder : %v", err)
			json.NewEncoder(w).Encode(map[string]interface{}{"msg": "edit fail : ", "data": doc})
		} else {
			doc, _ := client.Collection(db).Doc(fmt.Sprintf("%v", user.ID)).Get(ctx)
			json.NewEncoder(w).Encode(doc.Data())
		}
	case "DELETE":
		//delete data from id
		id := r.URL.Query()["id"]
		_, err := client.Collection(db).Doc(id[0]).Delete(ctx)
		if err != nil {
			log.Fatalf("An error has occurred : %v ", err)
			json.NewEncoder(w).Encode(map[string]interface{}{"msg": "delete fail"})
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"ID": id[0],
			})
		}

	}

}

//fetchData : get data all
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

//searchData use key data that be search, and target is speactifier column
func serchData(key string, target string, data []map[string]interface{}) (map[string]interface{}, error) {
	for i, v := range data {
		if key == v[target] {
			return data[i], nil
		}
	}
	return nil, errors.New("not Found data")
}
