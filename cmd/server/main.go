package main

import (
	"GolandRestApi/pkg/api/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	fmt.Println("Hello, World!")
	
	r := mux.NewRouter()

	// Define routes here
	r.HandleFunc("/example", handlers.ExampleHandler)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
