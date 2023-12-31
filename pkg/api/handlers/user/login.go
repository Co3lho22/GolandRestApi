package user

import (
	"GolandRestApi/pkg/config"
	"GolandRestApi/pkg/model"
	"GolandRestApi/pkg/repository"
	"GolandRestApi/pkg/service"
	"GolandRestApi/pkg/utils"
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
		service.HttpErrorResponse(logger,
			w,
			http.StatusBadRequest,
			"/login",
			"Invalid request format",
			err,
			utils.LogTypeError,
			"not able to get the username")
		return
	}

	var newUser *model.User
	newUser, err = repository.GetUserByUserName(logger, db, loginDetails.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.HttpErrorResponse(logger,
				w,
				http.StatusUnauthorized,
				"/login",
				"Invalid username or password",
				nil,
				utils.LogTypeWarn,
				loginDetails.Username)
			return
		}

		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/login",
			"Error creating user",
			err,
			utils.LogTypeError,
			loginDetails.Username)
		return
	}

	newUser.Username = loginDetails.Username
	err = service.CheckPasswordHash(logger, newUser, loginDetails.Password)
	if err != nil {
		service.HttpErrorResponse(logger,
			w,
			http.StatusUnauthorized,
			"/login",
			"Invalid username or password",
			nil,
			utils.LogTypeWarn,
			loginDetails.Username)
		return
	}

	var accessToken, refreshToken string
	accessToken, refreshToken, err = service.HandleTokensCreation(logger, db, cfg, loginDetails.Username)
	if err != nil {
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/login",
			"Server error handling the tokens",
			err,
			utils.LogTypeError,
			loginDetails.Username)
		return
	}

	var result bool
	result, err = repository.StoreRefreshTokenInDB(logger, db, refreshToken, loginDetails.Username)
	if err != nil {
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/login",
			"Server error storing refreshToken",
			err,
			utils.LogTypeError,
			loginDetails.Username)
		return
	} else if !result {
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/login",
			"Server failed to storing refreshToken",
			nil,
			utils.LogTypeWarn,
			loginDetails.Username)
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
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/login",
			"Error writing the response",
			err,
			utils.LogTypeError,
			loginDetails.Username)
		return
	}

	logger.WithField("username", loginDetails.Username).Info("User logged in with success")
	return
}


