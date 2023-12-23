package user

import (
	"GolandRestApi/pkg/model"
	"GolandRestApi/pkg/repository"
	"GolandRestApi/pkg/service"
	"GolandRestApi/pkg/utils"
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
		//logger.WithError(err).Error("Failed to deserialize the User object for the /register endpoint")
		//http.Error(w, "Invalid request format", http.StatusBadRequest)
		service.HttpErrorResponse(logger,
			w,
			http.StatusBadRequest,
			"/register",
			"Invalid request format",
			err,
			utils.LogTypeError,
			"not able to get the username")
		return
	}

	userExists, err := repository.UserExists(logger, db, newUser.Username, newUser.Email)
	if err != nil {
		//http.Error(w, "Error while verifying if user exists", http.StatusInternalServerError)
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/register",
			"Error while verifying if user exists",
			err,
			utils.LogTypeError,
			newUser.Username)
		return
	} else if userExists == true {
		//http.Error(w, "Username or Email already in use", http.StatusMethodNotAllowed)
		service.HttpErrorResponse(logger,
			w,
			http.StatusMethodNotAllowed,
			"/register",
			"Username or Email already in use",
			nil,
			utils.LogTypeInfo,
			newUser.Username)
		return
	}

	hashedPassword, err := service.HashPassword(logger, newUser)
	if err != nil {
		//http.Error(w, "Error hashing password", http.StatusInternalServerError)
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/register",
			"Error hashing password",
			err,
			utils.LogTypeError,
			newUser.Username)
		return
	}

	newUser.HashedPassword = hashedPassword
	err = repository.AddUser(logger, db, newUser, utils.UserRole)
	if err != nil {
		//http.Error(w, "Error adding user", http.StatusInternalServerError)
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/register",
			"Error adding user",
			err,
			utils.LogTypeError,
			newUser.Username)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("User successfully created"))
	if err != nil {
		//logger.WithError(err).Error("Error writing the response")
		//http.Error(w, "Error writing the response", http.StatusInternalServerError)
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/register",
			"Error writing response",
			err,
			utils.LogTypeError,
			newUser.Username)
		return
	}

	logger.WithField("username", newUser.Username).Info("User registered with success")
	return
}
