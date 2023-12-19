package handlers

import (
	"GolandRestApi/pkg/model"
	"GolandRestApi/pkg/repository"
	"GolandRestApi/pkg/service"
	"database/sql"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

// RegisterUser is an HTTP handler function for registering a new user.
// It takes a logrus.Logger instance for logging, a pointer to a sql.DB representing the database connection,
// a http.ResponseWriter for writing the HTTP response, and an http.Request for processing the HTTP request.
// This function expects a JSON-encoded User object in the request body and performs the following steps:
// 1. Deserialize the User object from the request body.
// 2. Check if a user with the same username or email already exists in the database.
// 3. If not, hash the user's password and add the user to the database.
// 4. Respond with appropriate HTTP status codes and messages.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to a sql.DB representing the database connection.
// w: An http.ResponseWriter for writing the HTTP response.
// r: An http.Request containing the HTTP request with a JSON-encoded User object in the request body.
func RegisterUser(logger *logrus.Logger, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var newUser model.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		logger.WithError(err).Error("Failed to deserialize the User object for the /register endpoint")
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	userExists, err := repository.UserExists(logger, db, newUser.Username, newUser.Email)
	if err != nil {
		logger.WithError(err).Error("Error verifying if user exists")
		http.Error(w, "Error while verifying if user exists", http.StatusInternalServerError)
		return
	} else if userExists == true {
		http.Error(w, "Username or Email already in use", http.StatusMethodNotAllowed)
		return
	}

	hashedPassword, err := service.HashPassword(logger, newUser)
	if err != nil {
		logger.WithError(err).Error("Error adding user")
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	newUser.HashedPassword = hashedPassword
	success, err := repository.AddUser(logger, db, newUser)
	if err != nil {
		logger.WithError(err).Error("Error adding user")
		http.Error(w, "Error adding user", http.StatusInternalServerError)
		return
	}

	if !success {
		http.Error(w, "Error while registering user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("User successfully created"))
	if err != nil {
		logger.WithError(err).Error("Error writing the response")
		http.Error(w, "Error writing the response", http.StatusInternalServerError)
		return
	}

	logger.WithField("username", newUser.Username).Info("User registered with success")
	return
}
