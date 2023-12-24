package repository

import (
	"GolandRestApi/pkg/model"
	"database/sql"
	"errors"
	"fmt"
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

// addUserWithoutRole adds a new user to the database with the provided user information and without assigning a role.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the sql.DB instance representing the database connection.
// user: A model.User struct containing the user information to be added.
//
// Returns an error if there is any issue while adding the user.
// If successful, the user is added to the database without errors.
func addUserWithoutRole(logger *logrus.Logger, db *sql.DB, user model.User) error {
	query := `INSERT INTO USERS (username, hashed_password, email, country, phone) VALUES (?, ?, ?, ?, ?)`
	result, err := db.Exec(query, user.Username, user.HashedPassword, user.Email, user.Country, user.Phone)
	if err != nil {
		logger.WithError(err).WithField("username", user.Username).Error("Error adding user without roles")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).WithField("username", user.Username).Error("Error getting rows affected")
		return err
	}

	if !(rowsAffected > 0) {
		logger.WithField("username", user.Username).Warn("No errors adding user, but no rows affected")
		return fmt.Errorf("no errors adding user %s, but no rows affected", user.Username)
	}

	logger.WithField("username", user.Username).Info("Add user without roles with success")
	return nil
}

// AddUser adds a new user to the database with the provided user information and assigns a specified role.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the sql.DB instance representing the database connection.
// user: A model.User struct containing the user information to be added.
// roleName: The name of the role to be assigned to the user.
//
// Returns an error if there is any issue while adding the user or setting the role.
// If successful, the user is added to the database with the specified role without errors.
func AddUser(logger *logrus.Logger, db *sql.DB, user model.User, roleName string) error {
	err := addUserWithoutRole(logger, db, user)
	if err != nil {
		return err
	}

	err = SetUserRole(logger, db, user.Username, roleName)
	if err != nil {
		return err
	}

	logger.WithField("username", user.Username).Info("Add user with success")
	return nil
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

// GetUserNameByUserId retrieves the username associated with a user ID from the database.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the SQL database instance.
// userId: The ID of the user whose username needs to be retrieved.
//
// Returns the username associated with the given user ID.
// Returns an error if the user is not found in the database or if there's any error during retrieval.
func GetUserNameByUserId(logger *logrus.Logger, db *sql.DB, userId int) (string, error) {
	query := `SELECT username FROM USERS WHERE id= ?`
	var username string
	err := db.QueryRow(query, userId).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithField("userId", userId).Info("Username not found in DB")
			return "", err
		}
		logger.WithError(err).WithField("userId", userId).Error("Error retrieving username from DB")
		return "", err
	}

	logger.WithField("userId", userId).Info("Get username by userId with success")
	return username, nil
}

// GetUserIdByUserName retrieves the user ID associated with a username from the database.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the SQL database instance.
// username: The username for which the user ID needs to be retrieved.
//
// Returns the user ID associated with the given username.
// Returns -1 and an error if the user is not found in the database or if there's any error during retrieval.
func GetUserIdByUserName(logger *logrus.Logger, db *sql.DB, username string) (int, error) {
	query := `SELECT id FROM USERS WHERE username= ?`
	var userId int
	err := db.QueryRow(query, userId).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithField("userName", username).Info("UserId not found in DB")
			return -1, err
		}
		logger.WithError(err).WithField("userName", username).Error("Error retrieving userId from DB")
		return -1, err
	}

	logger.WithField("userName", username).Info("Get userId by username with success")
	return userId, nil
}

// DeleteUser removes a user and associated records from the database based on the user ID.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the SQL database instance.
// userId: The ID of the user to be removed.
//
// Returns an error if there is any issue while deleting the user or associated records.
// If successful, the user and associated records are removed from the database without errors.
func DeleteUser(logger *logrus.Logger, db *sql.DB, userId int) error {
	tx, err := db.Begin()
	if err != nil {
		logger.WithError(err).WithField("userId", userId).Error("Error beginning the multiple queries")
		return err
	}

	if _, err := tx.Exec("DELETE FROM USER_ROLE WHERE user_id = ?", userId); err != nil {
		logger.WithError(err).WithField("userId", userId).Error("Error executing the first query to remove a user")
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM USER_AUTH WHERE user_id = ?", userId); err != nil {
		logger.WithError(err).WithField("userId", userId).Error("Error executing the second query to remove a user")
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM USERS WHERE id = ?", userId); err != nil {
		logger.WithError(err).WithField("userId", userId).Error("Error executing the third query to remove a user")
		tx.Rollback()
		return err
	}

	logger.WithField("userId", userId).Info("user removed successfully")
	return tx.Commit()
}
