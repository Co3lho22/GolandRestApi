package user

import (
	"GolandRestApi/pkg/repository"
	"database/sql"
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
		logger.Error("UserId not provided in the URL")
		http.Error(w, "UserId not provided", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		logger.WithError(err).Error("Invalid userId format")
		http.Error(w, "Invalid userId format", http.StatusBadRequest)
		return
	}

	success, err := repository.TokenRevocation(logger, db, userId)
	if err != nil || !success {
		http.Error(w, "Error during logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Logged out successfully"))
	if err != nil {
		logger.WithError(err).WithField("userId", userId).Error("Error writing response")
		return
	}

	logger.WithField("userId", userId).Info("User logged out successfully")
	return
}
