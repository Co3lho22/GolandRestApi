package handlers

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RefreshToken(logger *logrus.Logger, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	//TODO: Add the logic to verify refreshToken and return the new access, and refresh token if apply

	//1. Client sends the refresh token
	//2. Call the function VerifyToken to make sure the token in valid and not been tempered. If not error
	//3. Get the refresh token from the db and compare if they are the same. If not error
	//4. Generate a new access token and a new refresh token
	//4. Update/add the new refresh token to the db
	//5. Send the new access token and a new refresh token to the client

	//var refreshTokenDetails struct {
	//	RefreshToken string `json:"refreshToken"`
	//}

}
