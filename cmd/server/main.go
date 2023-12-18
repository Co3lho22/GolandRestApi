package main

import (
	"GolandRestApi/pkg/api/handlers"
	"GolandRestApi/pkg/config"
	"GolandRestApi/pkg/service"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func main() {
	cfg := config.NewConfig()
	serverPort := cfg.ServerPort

	logger, errLog := service.NewLogger(cfg)

	if errLog != nil {
		logger.Fatal("Could not initialize logger: ", errLog)
	}

	r := mux.NewRouter()

	logger.Info("Http server started on port ", serverPort, ".\n")
	
	// Define routes here
	r.HandleFunc("/login", handlers.LoginUser)
	r.HandleFunc("/register", handlers.RegisterUser)

	db, errDB := service.NewDBConnection(cfg)
	if errDB != nil {
		logger.Fatal("Could not connect to the database: ", errDB)
	}
	defer db.Close()

	errHttp := http.ListenAndServe(":"+strconv.Itoa(serverPort), r)
	if errHttp != nil {
		logger.Fatal("Error starting server: ", errHttp)
		return
	}
}
