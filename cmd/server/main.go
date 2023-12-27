package main

import (
	"GolandRestApi/pkg/api/handlers/admin"
	"GolandRestApi/pkg/api/handlers/middleware"
	"GolandRestApi/pkg/api/handlers/token"
	"GolandRestApi/pkg/api/handlers/user"
	"GolandRestApi/pkg/config"
	"GolandRestApi/pkg/service"
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

//TODO: 6. Dockerfile Usage
// The Dockerfile is used to containerize your application.
// This means you can run your application in a Docker container, which packages your application and its environment.
// This is useful for ensuring consistency across different environments and is beneficial for deployment and scaling,
// especially when used with orchestration tools like Kubernetes.

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

	r.Use(middleware.Authenticate(logger, db, cfg))

	// Routes
	mainRoutFormatted := "/api/" + cfg.APIVersion
	mainRoute := r.PathPrefix(mainRoutFormatted).Subrouter()

	// User routes
	userRoutes := mainRoute.PathPrefix("/user").Subrouter()
	userRoutes.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		user.LoginUser(logger, db, cfg, w, r)
	}).Methods("POST")
	userRoutes.HandleFunc("/logout/{userId}", func(w http.ResponseWriter, r *http.Request) {
		user.LogoutUser(logger, db, w, r)
	}).Methods("GET")
	userRoutes.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		user.RegisterUser(logger, db, w, r)
	}).Methods("POST")

	// Token routes
	tokenRoutes := mainRoute.PathPrefix("/token").Subrouter()
	tokenRoutes.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		token.Refresh(logger, db, cfg, w, r)
	}).Methods("POST")

	//// Admin routes
	adminRoutes := mainRoute.PathPrefix("/admin").Subrouter()
	adminRoutes.HandleFunc("/addUser", func(w http.ResponseWriter, r *http.Request) {
		admin.AddUser(logger, db, w, r)
	}).Methods("POST")
	adminRoutes.HandleFunc("/removeUser/{userId}", func(w http.ResponseWriter, r *http.Request) {
		admin.RemoveUser(logger, db, w, r)
	}).Methods("DELETE")

	err = http.ListenAndServe(":"+strconv.Itoa(serverPort), r)
	if err != nil {
		logger.Fatal("Error starting server: ", err)
		return
	}
}
