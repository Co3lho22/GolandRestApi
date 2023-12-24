package admin

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

// AddUser handles the creation of a new user by an administrator.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the SQL database connection.
// w: The http.ResponseWriter to write the response to.
// r: The HTTP request containing user and role information in JSON format.
//
// This function expects a JSON request body with user details (model.User) and a roleName.
// It first validates the request format and checks if the provided username and email are unique.
// If the user details are valid and unique, the function hashes the password, creates the user with the specified role,
// and sends a success response. If any error occurs during the process, an appropriate error response is sent.
func AddUser(logger *logrus.Logger, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var AddUserDetails struct {
		User     model.User `json:"user"`
		RoleName string     `json:"roleName"`
	}

	err := json.NewDecoder(r.Body).Decode(&AddUserDetails)
	if err != nil {
		service.HttpErrorResponse(logger,
			w,
			http.StatusBadRequest,
			"/admin/adduser",
			"Invalid request format",
			err,
			utils.LogTypeError,
			"not able to get the username")
		return
	}

	userExists, err := repository.UserExists(logger, db, AddUserDetails.User.Username, AddUserDetails.User.Email)
	if err != nil {
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/admin/adduser",
			"Error while verifying if user exists",
			err,
			utils.LogTypeError,
			AddUserDetails.User.Username)
		return
	} else if userExists == true {
		service.HttpErrorResponse(logger,
			w,
			http.StatusMethodNotAllowed,
			"/admin/adduser",
			"Username or Email already in use",
			nil,
			utils.LogTypeInfo,
			AddUserDetails.User.Username)
		return
	}

	hashedPassword, err := service.HashPassword(logger, AddUserDetails.User)
	if err != nil {
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/admin/adduser",
			"Error hashing password",
			err,
			utils.LogTypeError,
			AddUserDetails.User.Username)
		return
	}

	AddUserDetails.User.HashedPassword = hashedPassword
	err = repository.AddUser(logger, db, AddUserDetails.User, AddUserDetails.RoleName)
	if err != nil {
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/admin/adduser",
			"Error creating user via admin",
			err,
			utils.LogTypeError,
			AddUserDetails.User.Username)
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
			"/admin/addUser",
			"Error writing the response",
			err,
			utils.LogTypeError,
			AddUserDetails.User.Username)
		return
	}

	logger.WithField("username", AddUserDetails.User.Username).Info("User created with success")
	return
}
