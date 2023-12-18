package main

import (
	"GolandRestApi/pkg/api/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	// Define routes here
	r.HandleFunc("/example", handlers.ExampleHandler)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
