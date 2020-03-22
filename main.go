package main

import (
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/user", UserHandler)
	http.ListenAndServe(getPort(), nil)
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
