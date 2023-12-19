package repository

import (
	"GolandRestApi/pkg/model"
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
)

func UserExists(logger *logrus.Logger, db *sql.DB, username string, email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM USERS WHERE username = ? OR email = ?)`
	err := db.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": username,
		}).WithError(err).Error("Error checking if user exists")

		return false, errors.New("error checking if user exists")
	}
	return count >= 1, nil
}

func AddUser(logger *logrus.Logger, db *sql.DB, user model.User) (bool, error) {
	query := `INSERT INTO USERS (username, hashed_password, email, country, phone) VALUES (?, ?, ?, ?, ?)`
	result, err := db.Exec(query, user.Username, user.HashedPassword, user.Email, user.Country, user.Phone)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": user.Username,
		}).WithError(err).Error("Error adding user")

		return false, errors.New("error adding user")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Errorf("Error getting rows affected: %v", err)
		return false, err
	}
	return rowsAffected > 0, nil
}
