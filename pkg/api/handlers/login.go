package handlers

import (
	"GolandRestApi/pkg/model"
	"GolandRestApi/pkg/repository"
	"GolandRestApi/pkg/service"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

// LoginUser is an HTTP handler function for user login.
// It takes a logrus.Logger instance for logging, a pointer to a sql.DB representing the database connection,
// a http.ResponseWriter for writing the HTTP response, and an http.Request for processing the HTTP request.
// This function expects a JSON-encoded structure containing the username and password in the request body.
// It performs the following steps:
// 1. Deserialize the username and password from the request body.
// 2. Retrieve the user from the database by username.
// 3. If the user doesn't exist or the password is invalid, return an appropriate error response.
// 4. If the username and password are valid, respond with an HTTP status indicating successful login.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to a sql.DB representing the database connection.
// w: An http.ResponseWriter for writing the HTTP response.
// r: An http.Request containing the HTTP request with JSON-encoded login details.
func LoginUser(logger *logrus.Logger, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	var loginDetails struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginDetails)
	if err != nil {
		logger.WithError(err).Error("Failed to deserialize the User object for the /login endpoint")
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	var newUser *model.User
	newUser, err = repository.GetUserByUserName(logger, db, loginDetails.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		logger.WithError(err).
			WithField("username", loginDetails.Username).
			Error("Error retrieving user from DB")
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	newUser.Username = loginDetails.Username
	if err := service.CheckPasswordHash(logger, newUser, loginDetails.Password); err != nil {
		logger.
			WithField("username", loginDetails.Username).
			Warn("Invalid password")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	responseMessage := "User successfully logged in"
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(responseMessage))
	if err != nil {
		logger.WithError(err).Error("Error writing the response")
		http.Error(w, "Error writing the response", http.StatusInternalServerError)
		return
	}

	logger.WithField("username", loginDetails.Username).Info("User logged in with success")
	return
}
