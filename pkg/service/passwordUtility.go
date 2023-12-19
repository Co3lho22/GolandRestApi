package service

import (
	"GolandRestApi/pkg/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

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
