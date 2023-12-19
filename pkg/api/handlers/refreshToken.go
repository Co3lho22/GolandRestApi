package handlers

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RefreshToken(logger *logrus.Logger, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	//Expects a json body:
	//{
	//	"refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	//}
	//
	//1. Client sends the refresh token
	//2. Verify the refresh token
	//3. Generate a new access token and a new refresh token
	//4. Update/add the new refresh token to the db
	//5. Send the new access token and a new refresh token to the client

	//var refreshTokenDetails struct {
	//	RefreshToken string `json:"refreshToken"`
	//}
}
