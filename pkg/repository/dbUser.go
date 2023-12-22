package repository

import (
	"GolandRestApi/pkg/model"
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
)

// UserExists checks if a user with the given username or email already exists in the database.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the sql.DB instance representing the database connection.
// username: The username to be checked for existence.
// email: The email to be checked for existence.
//
// Returns true if a user with the provided username or email exists in the database, false otherwise.
// Returns an error if there is an issue with the database query.
func UserExists(logger *logrus.Logger, db *sql.DB, username string, email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM USERS WHERE username = ? OR email = ?`
	err := db.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		logger.WithError(err).WithField("username", username).Error("Error checking if user exists")
		return false, err
	}

	logger.WithField("username", username).Info("Verify if user exists with success")
	return count >= 1, nil
}

// AddUser adds a new user to the database with the provided user information.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the sql.DB instance representing the database connection.
// user: A model.User struct containing the user information to be added.
//
// Returns true if the user is successfully added to the database, false otherwise.
// Returns an error if there is an issue with the database query.
func AddUser(logger *logrus.Logger, db *sql.DB, user model.User) (bool, error) {
	query := `INSERT INTO USERS (username, hashed_password, email, country, phone) VALUES (?, ?, ?, ?, ?)`
	result, err := db.Exec(query, user.Username, user.HashedPassword, user.Email, user.Country, user.Phone)
	if err != nil {
		logger.WithError(err).WithField("username", user.Username).Error("Error adding user")
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).WithField("username", user.Username).Error("Error getting rows affected")
		return false, err
	}

	if !(rowsAffected > 0) {
		logger.WithField("username", user.Username).Warn("No errors adding user, but no rows affected")
		return false, nil
	}

	logger.WithField("username", user.Username).Info("Add user with success")
	return true, nil
}

// GetUserByUserName retrieves user information from the database based on the provided username.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the sql.DB instance representing the database connection.
// username: The username for which user information should be retrieved.
//
// Returns a pointer to a model.User struct containing the user information if found in the database.
// Returns sql.ErrNoRows if no user with the provided username is found.
// Returns an error if there is an issue with the database query.
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
		logger.WithError(err).WithField("username", username).Error("Error retrieving user from DB")
		return nil, err
	}

	logger.WithField("username", username).Info("Get user by username with success")
	return &user, nil
}
