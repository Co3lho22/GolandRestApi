package service

import (
	"GolandRestApi/pkg/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hashed password for a given user's plaintext password.
// It takes a logrus.Logger instance for logging, and a model.User struct representing the user.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// user: A model.User struct containing the user's information, including the plaintext password.
//
// Returns the bcrypt hashed password as a string and an error, if any. If there is an error
// during password hashing, it returns an empty string and the error.
func HashPassword(logger *logrus.Logger, user model.User) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": user.Username,
		}).WithError(err).Error("Error adding user")
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a given plaintext password with a bcrypt hashed password.
// It takes a logrus.Logger instance for logging, a pointer to a model.User struct
// representing the user, and the plaintext password to check.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// user: A pointer to a model.User struct containing the user's information and hashed password.
// password: The plaintext password to check against the hashed password.
//
// Returns an error if the password does not match the hashed password, and logs an error message
// with the username if the comparison fails. Returns nil if the passwords match.
func CheckPasswordHash(logger *logrus.Logger, user *model.User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": user.Username,
		}).WithError(err).Error("Error checking passwordHash")
		return err
	}
	return nil
}
