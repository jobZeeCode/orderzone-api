package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type orderDetailStruct struct {
	ID interface{}
	OrderID interface{}
	MenuID interface{}
	Amount interface{}
	NetPrice interface{}
}

//OrderDetailHandler : manage shop database
func OrderDetailHandler (w http.ResponseWriter, r *http.Request) {
	client, ctx = Conect()
	db := "order-detail"
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
			data, err := serchData(id[0], "ID", fetchData(db))
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"msg": "not found"})
			} else {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(data)
			}
		} else {
			//Not have query
			data := fetchData(db)
			if len(data) > 0 {
				//
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(fetchData(db))
			} else {
				json.NewEncoder(w).Encode(map[string]interface{}{"msg": "not found"})
			}
		}
		break
	case "POST":
		//insert data
		var orderDetail orderDetailStruct 
		err := json.NewDecoder(r.Body).Decode(&orderDetail)
		if err != nil {
			log.Fatalf("Fail Decoder : %v ", err)
		}
		ref := client.Collection(db).NewDoc()
		orderDetail.ID = ref.ID
		doc, err := ref.Set(ctx, orderDetail)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"msg": "add fail : ", "data": doc})
		} else {
			json.NewEncoder(w).Encode(orderDetail)
		}
		break
	case "PUT":
		//Edit data
		var orderDetail orderDetailStruct 
		err := json.NewDecoder(r.Body).Decode(&orderDetail)
		if err != nil {
			log.Fatalf("Failt Decoder : %v ", err)
		}
		doc, err := client.Collection(db).Doc(fmt.Sprintf("%v", orderDetail.ID)).Set(ctx, orderDetail)
		if err != nil {
			log.Fatalf("Fail Decoder : %v", err)
			json.NewEncoder(w).Encode(map[string]interface{}{"msg": "edit fail : ", "data": doc})
		} else {
			doc, _ := client.Collection(db).Doc(fmt.Sprintf("%v", orderDetail.ID)).Get(ctx)
			json.NewEncoder(w).Encode(doc.Data())
		}
		break
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
		break
	default:
		break
	}
}