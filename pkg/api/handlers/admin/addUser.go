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

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("User successfully created"))
	if err != nil {
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/admin/adduser",
			"Error writing response",
			err,
			utils.LogTypeError,
			AddUserDetails.User.Username)
		return
	}

	logger.WithField("username", AddUserDetails.User.Username).Info("User created with success")
	return
}
