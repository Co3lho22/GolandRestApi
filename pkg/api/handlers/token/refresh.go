package token

import (
	"GolandRestApi/pkg/config"
	"GolandRestApi/pkg/repository"
	"GolandRestApi/pkg/service"
	"GolandRestApi/pkg/utils"
	"database/sql"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Refresh handles the refresh of access tokens for an authenticated user.
// It receives a refresh token in the request, verifies its authenticity and validity,
// and generates a new access token and refresh token pair if the provided token is valid.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the SQL database connection.
// cfg: A pointer to the config.Config struct which contains JWT configuration.
// w: The HTTP response writer to send the response.
// r: The HTTP request to process.
//
// If the refresh token is valid and matches the one stored in the database for the user,
// it returns a new access token and refresh token in the response along with a status code 200 (OK).
// If any errors occur during token verification, generation, or storage, it returns an appropriate
// HTTP error response with the corresponding status code and error message.
func Refresh(logger *logrus.Logger, db *sql.DB, cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var refreshDetails struct {
		RefreshToken string `json:"refreshToken"`
	}

	err := json.NewDecoder(r.Body).Decode(&refreshDetails)
	if err != nil {
		//logger.WithError(err).Error("Failed to deserialize the User object for the /refresh endpoint")
		//http.Error(w, "Invalid request format", http.StatusBadRequest)
		service.HttpErrorResponse(logger,
			w,
			http.StatusBadRequest,
			"/refresh",
			"Invalid request format",
			err,
			utils.LogTypeError,
			"not able to get the username")
		return
	}

	userName, _, err := service.ExtractClaimsFromToken(logger, refreshDetails.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		service.HttpErrorResponse(logger,
			w,
			http.StatusBadRequest,
			"/refresh",
			"Invalid token",
			err,
			utils.LogTypeError,
			"not able to get the username")
		return
	}

	err = service.VerifyToken(logger, cfg, userName, refreshDetails.RefreshToken)
	if err != nil {
		//http.Error(w, "Invalid token", http.StatusBadRequest)
		service.HttpErrorResponse(logger,
			w,
			http.StatusBadRequest,
			"/refresh",
			"Invalid token",
			err,
			utils.LogTypeError,
			userName)
		return
	}

	dbRefreshToken, err := repository.RetrieveRefreshTokenFromDB(logger, db, userName)
	if err != nil {
		//http.Error(w, "Server error retrieving refreshToken from DB", http.StatusInternalServerError)
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/refresh",
			"Server error retrieving refreshToken from DB",
			err,
			utils.LogTypeError,
			userName)
		return
	}

	if dbRefreshToken != refreshDetails.RefreshToken {
		//http.Error(w, "Invalid token", http.StatusBadRequest)
		service.HttpErrorResponse(logger,
			w,
			http.StatusBadRequest,
			"/refresh",
			"Invalid token",
			nil,
			utils.LogTypeWarn,
			userName)
		return
	}

	var newAccessToken, newRefreshToken string
	newAccessToken, newRefreshToken, err = service.HandleTokensCreation(logger, db, cfg, userName)
	if err != nil {
		//logger.WithError(err).Error("Error creating the token for the /refresh")
		//http.Error(w, "Server error handling the tokens", http.StatusInternalServerError)
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/refresh",
			"Server error handling the tokens",
			err,
			utils.LogTypeError,
			userName)
		return
	}

	var result bool
	result, err = repository.StoreRefreshTokenInDB(logger, db, newRefreshToken, userName)
	if err != nil {
		http.Error(w, "Server error storing refreshToken", http.StatusInternalServerError)
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/refresh",
			"Server error storing refreshToken",
			err,
			utils.LogTypeError,
			userName)
		return
	} else if !result {
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/refresh",
			"Server failed to storing refreshToken",
			nil,
			utils.LogTypeWarn,
			userName)
		return
	}

	response := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		//logger.WithError(err).Error("Error writing the response")
		//http.Error(w, "Error writing the response", http.StatusInternalServerError)
		service.HttpErrorResponse(logger,
			w,
			http.StatusInternalServerError,
			"/login",
			"Error writing the response",
			err,
			utils.LogTypeError,
			userName)
		return
	}

}
