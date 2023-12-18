package main

import (
	"GolandRestApi/pkg/api/handlers"
	"GolandRestApi/pkg/config"
	"GolandRestApi/pkg/service"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {
	cfg := config.NewConfig()
	serverPort := cfg.ServerPort

	logger, errLog := service.NewLogger(cfg)

	if errLog != nil {
		log.Fatalf("Could not initialize logger: %v", errLog)
	}

	logger.Info("Application started")

	fmt.Printf("Http server started on port %d\n", serverPort)

	r := mux.NewRouter()

	// Define routes here
	r.HandleFunc("/login", handlers.LoginUser)
	r.HandleFunc("/register", handlers.RegisterUser)

	db, errDB := service.NewDBConnection(cfg)
	if errDB != nil {
		log.Fatalf("Could not connect to the database: %v", errDB)
	}
	defer db.Close()

	err := http.ListenAndServe(":"+strconv.Itoa(serverPort), r)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
