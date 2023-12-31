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

	message := "User successfully created"
	response := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/user/register",
			"Error writing the response",
			err,
			utils.LogTypeError,
			newUser.Username)
		return
	}

	logger.WithField("username", newUser.Username).Info("User registered with success")
	return
}
