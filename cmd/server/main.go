package main

import (
	"GolandRestApi/pkg/api/handlers"
	"GolandRestApi/pkg/config"
	"GolandRestApi/pkg/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {
	cfg := config.NewConfig()
	serverPort := cfg.ServerPort

	logger, err := service.NewLogger(cfg)

	if err != nil {
		if logger != nil {
			logger.WithError(err).Fatal("Could not initialize logger")
		} else {
			log.Fatalf("Could not initialize logger: %v", err)
		}
	}

	r := mux.NewRouter()

	logger.Info("Http server started on port ", serverPort, ".")

	db, err := service.NewDBConnection(logger, cfg)
	if err != nil {
		logger.Fatal("Could not connect to the database: ", err)
	}
	defer db.Close()

	// Define routes here
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginUser(logger, db, w, r)
	})
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterUser(logger, db, w, r)
	})

	err = http.ListenAndServe(":"+strconv.Itoa(serverPort), r)
	if err != nil {
		logger.Fatal("Error starting server: ", err)
		return
	}
}
