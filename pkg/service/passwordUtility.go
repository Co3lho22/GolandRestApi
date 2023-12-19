package service

import (
	"GolandRestApi/pkg/model"
	"errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(logger *logrus.Logger, user model.User) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": user.Username,
		}).WithError(err).Error("Error adding user")
		return "", errors.New("error hashing password")
	}
	return string(bytes), nil
}

func CheckPasswordHash(logger *logrus.Logger, user model.User, password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": user.Username,
		}).WithError(err).Error("Error checking passwordHash")
		return false, errors.New("error checking passwordHash")
	}
	return true, nil
}
