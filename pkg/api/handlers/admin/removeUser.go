package admin

import (
	"GolandRestApi/pkg/repository"
	"GolandRestApi/pkg/service"
	"GolandRestApi/pkg/utils"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func RemoveUser(logger *logrus.Logger, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	//TODO: what i need to remove?
	// 1. USERS
	// 2. USER_AUTH
	// 3. USER_ROLE
	vars := mux.Vars(r)
	userIdStr, ok := vars["userId"]
	if !ok {
		service.HttpErrorResponse(logger,
			w,
			http.StatusBadRequest,
			"/admin/removeUser",
			"User ID is required",
			nil, utils.LogTypeWarn,
			"")
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		service.HttpErrorResponse(logger,
			w,
			http.StatusBadRequest,
			"/admin/removeUser",
			"Invalid User ID format",
			err,
			utils.LogTypeError,
			"")
		return
	}

	username, err := repository.GetUserNameByUserId(logger, db, userId)
	if err != nil {
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/admin/removeUser",
			"Error retrieving user",
			err,
			utils.LogTypeError,
			"")
		return
	}

	if username == "" {
		service.HttpErrorResponse(logger,
			w,
			http.StatusNotFound,
			"/admin/removeUser",
			"User not found",
			nil, utils.LogTypeWarn,
			"")
		return
	}

	if err := repository.DeleteUser(logger, db, userId); err != nil {
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/admin/removeUser",
			"Error deleting user",
			err,
			utils.LogTypeError,
			username)
		return
	}

	message := "User successfully removed"
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
			"/admin/removeUser",
			"Error writing the response",
			err,
			utils.LogTypeError,
			username)
		return
	}

	logger.WithField("userId", userId).Info("User removed successfully")
	return
}
