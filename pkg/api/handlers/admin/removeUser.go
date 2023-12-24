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

// RemoveUser handles the removal of a user by an administrator.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the SQL database connection.
// w: The http.ResponseWriter to write the response to.
// r: The HTTP request containing the user ID as a path variable.
//
// This function retrieves the user ID from the request, validates it, and attempts to delete the user from the database.
// It performs the following steps:
// 1. Parse the user ID from the path variable.
// 2. Retrieve the username associated with the user ID from the database.
// 3. Check if the user exists; if not, return a not found response.
// 4. Attempt to delete the user from the database.
// 5. Send a success response if the user is successfully removed or an error response if any issues occur.
//
// Note: This function deletes records from multiple database tables (USERS, USER_AUTH, USER_ROLE) associated with the user.
func RemoveUser(logger *logrus.Logger, db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
