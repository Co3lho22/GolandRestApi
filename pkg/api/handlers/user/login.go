package user

import (
	"GolandRestApi/pkg/config"
	"GolandRestApi/pkg/model"
	"GolandRestApi/pkg/repository"
	"GolandRestApi/pkg/service"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

// LoginUser handles user authentication by verifying the provided username and password.
// Upon successful authentication, it generates and returns an access token and a refresh token.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the sql.DB instance representing the database connection.
// cfg: A pointer to the config.Config struct which contains JWT configuration details.
// w: The http.ResponseWriter to write the HTTP response.
// r: The http.Request representing the HTTP request with user login details in JSON format.
//
// Responds with a JSON object containing the access token and refresh token upon successful login.
// If login details are invalid, it returns an error response with an appropriate HTTP status code.
func LoginUser(logger *logrus.Logger, db *sql.DB, cfg *config.Config, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

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
			logger.WithField("username", loginDetails.Username).Warn("Invalid username or password")
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	newUser.Username = loginDetails.Username
	err = service.CheckPasswordHash(logger, newUser, loginDetails.Password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	var accessToken, refreshToken string
	accessToken, refreshToken, err = service.HandleTokensCreation(logger, db, cfg, loginDetails.Username)
	if err != nil {
		logger.WithError(err).Error("Error creating the token for the /login")
		http.Error(w, "Server error handling the tokens", http.StatusInternalServerError)
		return
	}

	response := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.WithError(err).Error("Error writing the response")
		http.Error(w, "Error writing the response", http.StatusInternalServerError)
		return
	}

	logger.WithField("username", loginDetails.Username).Info("User logged in with success")
	return
}
