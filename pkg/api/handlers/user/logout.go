package user

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

// LogoutUser handles the user logout process, which involves revoking the user's access and refresh tokens.
// It extracts the user ID from the URL parameters, revokes the user's tokens, and sends a response indicating
// successful logout.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the SQL database connection.
// w: The HTTP response writer to send the response.
// r: The HTTP request to process.
//
// If the user ID is not provided in the URL or has an invalid format, it returns an appropriate HTTP error response
// with a status code 400 (Bad Request). If token revocation encounters an error or is not successful, it returns
// an HTTP error response with a status code 500 (Internal Server Error). Upon successful logout, it sends an HTTP
// response with a status code 200 (OK) indicating that the user has logged out.
func LogoutUser(logger *logrus.Logger, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr, ok := vars["userId"]

	if !ok {
		//logger.Error("UserId not provided in the URL")
		//http.Error(w, "UserId not provided", http.StatusBadRequest)
		service.HttpErrorResponse(logger,
			w,
			http.StatusBadRequest,
			"/logout",
			"UserId not provided",
			nil,
			utils.LogTypeWarn,
			"not able to get the username")
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		//logger.WithError(err).Error("Invalid userId format")
		//http.Error(w, "Invalid userId format", http.StatusBadRequest)
		service.HttpErrorResponse(logger,
			w,
			http.StatusBadRequest,
			"/logout",
			"Invalid userId format",
			err,
			utils.LogTypeError,
			"not able to get the username")
		return
	}

	username, err := repository.GetUserNameByUserId(logger, db, userId)
	if err != nil {
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/logout",
			"Error getting userName by userId",
			err,
			utils.LogTypeError,
			"with userId: "+strconv.Itoa(userId))
		return
	}

	success, err := repository.TokenRevocation(logger, db, userId)
	if err != nil {
		//http.Error(w, "Error during logout", http.StatusInternalServerError)
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/logout",
			"Error revoking token",
			err,
			utils.LogTypeError,
			username)
		return
	} else if !success {
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/logout",
			"Failed revoking token",
			nil,
			utils.LogTypeWarn,
			username)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Logged out successfully"))
	if err != nil {
		//logger.WithError(err).WithField("userId", userId).Error("Error writing response")
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/logout",
			"Error writing response",
			err,
			utils.LogTypeError,
			username)
		return
	}

	message := "Logged out successfully"
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
			"/user/logout",
			"Error writing the response",
			err,
			utils.LogTypeError,
			username)
		return
	}

	logger.WithField("userId", userId).Info("User logged out successfully")
	return
}
