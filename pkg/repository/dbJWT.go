package repository

import (
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
)

// RetrieveRefreshTokenFromDB retrieves the refresh token associated with a user from the database.
// It queries the database for the refresh token based on the user's ID.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the sql.DB instance representing the database connection.
// userId: The ID of the user for whom the refresh token should be retrieved.
//
// Returns the retrieved refresh token as a string and an error.
// If the refresh token is not found, it returns an empty string and sql.ErrNoRows.
func RetrieveRefreshTokenFromDB(logger *logrus.Logger, db *sql.DB, userId int) (string, error) {
	query := "SELECT refresh_token FROM USER_AUTH WHERE user_id = ?"
	var refreshToken string
	err := db.QueryRow(query, userId).Scan(&refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithField("userId", userId).Info("Refresh token for this user not found in DB")
			return "", err
		}
		logger.WithError(err).WithField("userId", userId).Error("Error retrieving refreshToken from DB")
		return "", err
	}

	logger.WithField("userId", userId).Info("RefreshToken retrieve with success")
	return refreshToken, err
}

// StoreRefreshTokenInDB stores a refresh token in the database for a user.
// It inserts or updates the refresh token associated with the user's ID in the USER_AUTH table.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the sql.DB instance representing the database connection.
// refreshToken: The refresh token to be stored.
// userId: The ID of the user for whom the refresh token should be stored.
//
// Returns a boolean indicating whether the operation was successful and an error, if any.
func StoreRefreshTokenInDB(logger *logrus.Logger, db *sql.DB, refreshToken string, userName string) (bool, error) {
	query := "INSERT INTO USER_AUTH (user_id, refresh_token) VALUES ((SELECT id FROM USERS WHERE username = ?), ?) ON DUPLICATE KEY UPDATE refresh_token = ?"
	result, err := db.Exec(query, userName, refreshToken)
	if err != nil {
		logger.WithError(err).WithField("username", userName).Error("Error storing refreshToken")
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).WithField("username", userName).
			Error("Error getting the number of rows affected when trying to store the refreshToken")
		return false, err
	}

	if !(rowsAffected > 0) {
		logger.WithError(err).WithField("username", userName).
			Warn("No errors storing refreshToken in DB, but rows affected <= 0")
		return false, nil
	}

	logger.WithField("username", userName).Info("RefreshToken stored with success")
	return true, nil
}

// TokenRevocation revokes a user's refresh token by setting it to NULL in the database.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the sql.DB instance representing the database connection.
// userId: The ID of the user for whom the refresh token should be revoked.
//
// Returns a boolean indicating whether the operation was successful and an error, if any.
func TokenRevocation(logger *logrus.Logger, db *sql.DB, userId int) (bool, error) {
	query := "UPDATE USER_AUTH SET refresh_token = NULL WHERE user_id = ?;"
	result, err := db.Exec(query, userId)
	if err != nil {
		logger.WithError(err).WithField("userId", userId).Error("Error revoking token")
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).WithField("userId", userId).
			Error("Error getting the number of rows affected when trying to store the revoking a token")
		return false, err
	}

	if !(rowsAffected > 0) {
		logger.WithError(err).WithField("userId", userId).
			Warn("No errors storing revoking token, but rows affected <= 0")
		return false, nil
	}

	logger.WithField("userId", userId).Info("Token revoked with success")
	return true, nil
}
