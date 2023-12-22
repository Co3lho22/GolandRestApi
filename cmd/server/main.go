package main

import (
	"GolandRestApi/pkg/api/handlers"
	"GolandRestApi/pkg/config"
	"GolandRestApi/pkg/service"
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// main is the entry point of the GoLandRestApi application. It initializes and configures the HTTP server,
// sets up database connection, and defines the API routes for login and user registration.
//
// The main function uses the gorilla/mux router for handling HTTP requests.
// It also initializes a logger and database connection based on the provided configuration.
// The server listens on the specified port and handles incoming HTTP requests.
//
// It logs initialization errors and server start status.
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
		logger.WithError(err).Warn("Could not connect to the database")
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.WithError(err).Fatal("Could not close db")
		}
	}(db)

	// Define routes here
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginUser(logger, db, cfg, w, r)
	})
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		handlers.LogoutUser(logger, db, w, r)
	})
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterUser(logger, db, w, r)
	})
	r.HandleFunc("/refreshToken", func(w http.ResponseWriter, r *http.Request) {
		handlers.RefreshToken(logger, db, w, r)
	})

	err = http.ListenAndServe(":"+strconv.Itoa(serverPort), r)
	if err != nil {
		logger.Fatal("Error starting server: ", err)
		return
	}
}
