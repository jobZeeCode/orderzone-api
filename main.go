package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello World google App Engine with Golang")
	})
	http.HandleFunc("/user", UserHandler)
	http.HandleFunc("/user/login", PasswordHandler)
	http.HandleFunc("/shop", ShopHandler)
	http.HandleFunc("/menu", MenuHandler)
	http.HandleFunc("/order", OrderHandler)
	http.HandleFunc("/order/detail", OrderDetailHandler)
	http.ListenAndServe(getPort(), nil)
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
