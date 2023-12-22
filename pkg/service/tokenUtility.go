package service

import (
	"GolandRestApi/pkg/config"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"time"
)

// CreateToken generates a JSON Web Token (JWT) with the provided username and expiration time.
// It uses the secret key from the application configuration to sign the token. The expiration time
// is specified as a duration string (e.g., "15m" for 15 minutes). Note: This function is used to
// create the access and refresh token, just use the env variables for them.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// cfg: A pointer to the config.Config struct which contains the JWT secret key.
// username: The username to be included in the JWT claims.
// expirationTime: The duration for which the JWT will be valid.
//
// Returns the JWT token as a string and an error, if any. If there is an error during token creation,
// it returns an empty string and the error.
func CreateToken(logger *logrus.Logger, cfg *config.Config, username string, expirationTime string) (string, error) {
	var secretKey = []byte(cfg.JWTSecretKey)
	expirationDuration, err := time.ParseDuration(expirationTime)
	if err != nil {
		logger.WithField("username", username).
			WithError(err).
			Error("Invalid JWT expiration time format")
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(expirationDuration).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		logger.WithError(err).
			WithField("username", username).
			Error("Error creating the JWT token")
		return "", err
	}

	logger.WithField("username", username).
		Infof("JWT token generated")
	return tokenString, nil
}

// VerifyToken verifies the authenticity and validity of a JSON Web Token (JWT) using the provided secret key.
// It checks if the token is correctly signed and has not expired.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// cfg: A pointer to the config.Config struct which contains the JWT secret key.
// username: The username associated with the token (used for logging purposes).
// tokenString: The JWT token to be verified.
//
// Returns an error if the token is invalid, expired, or if there's any error during verification.
// Returns nil if the token is valid.
func VerifyToken(logger *logrus.Logger, cfg *config.Config, username string, tokenString string) error {
	var secretKey = []byte(cfg.JWTSecretKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		logger.WithError(err).
			WithField("username", username).
			Error("Error parsing the token")
		return err
	}

	if !token.Valid {
		logger.WithField("username", username).
			Warnf("Invalid token: %s", tokenString)
		return errors.New("invalid token")
	}

	return nil
}
