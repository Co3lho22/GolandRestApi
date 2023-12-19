package repository

import (
	"GolandRestApi/pkg/model"
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
)

func UserExists(logger *logrus.Logger, db *sql.DB, username string, email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM USERS WHERE username = ? OR email = ?`
	err := db.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": username,
		}).WithError(err).Error("Error checking if user exists")

		return false, err
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
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": user.Username,
		}).WithError(err).Error("Error getting rows affected")
		return false, err
	}
	return rowsAffected > 0, nil
}

func GetUserByUserName(logger *logrus.Logger, db *sql.DB, username string) (*model.User, error) {
	query := `SELECT id, hashed_password, email, country, phone FROM USERS WHERE username= ?`
	var user model.User
	user.Username = username
	err := db.QueryRow(query, username).Scan(&user.ID, &user.HashedPassword, &user.Email, &user.Country, &user.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithField("username", username).Info("User not found in DB")
			return nil, err
		}
		logger.WithError(err).
			WithField("username", username).
			Error("Error retrieving user from DB")
		return nil, err
	}

	return &user, nil
}
