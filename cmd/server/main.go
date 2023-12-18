package main

import (
	"GolandRestApi/pkg/api/handlers"
	"GolandRestApi/pkg/config"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func main() {
	cfg := config.NewConfig()
	serverPort := cfg.ServerPort

	fmt.Printf("Http server started on port %d\n", serverPort)

	r := mux.NewRouter()

	// Define routes here
	r.HandleFunc("/login", handlers.LoginUser)
	r.HandleFunc("/register", handlers.RegisterUser)

	err := http.ListenAndServe(":"+strconv.Itoa(serverPort), r)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
