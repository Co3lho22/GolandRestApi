package handlers

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"net/http"
)

func LogoutUser(logger *logrus.Logger, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	//TODO: Call the function TokenRevocation to remove the refresh token form the db
}
