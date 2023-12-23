package service

import (
	"GolandRestApi/pkg/config"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"time"
)

// createToken generates a JSON Web Token (JWT) with the provided username and expiration time.
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
func createToken(logger *logrus.Logger, cfg *config.Config, username string, expirationTime string) (string, error) {
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
		logger.WithError(err).WithField("username", username).Error("Error parsing the token")
		return err
	}

	if !token.Valid {
		logger.WithField("username", username).
			Warnf("Invalid token: %s", tokenString)
		return errors.New("invalid token")
	}

	return nil
}

// extractMapClaimsFromToken extracts the MapClaims from a JWT token without verifying its signature.
// It parses the provided token string and retrieves the claims as a map.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// tokenString: The JWT token to be parsed.
//
// Returns a jwt.MapClaims containing the claims from the token, or an error if parsing fails.
func extractMapClaimsFromToken(logger *logrus.Logger, tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		logger.WithError(err).WithField("token", tokenString).Error("Error parsing unverified token")
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.WithField("token", tokenString).Warn("Unable to retrieve claims from token")
		return nil, errors.New("invalid token claims")
	}

	logger.WithField("token", tokenString).Info("Token mapClaims retrieved successfully")
	return claims, nil
}

// ExtractClaimsFromToken extracts the username and expiration (exp) claims from a JWT token.
// It parses the provided token string, retrieves the claims, and returns them as strings.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// tokenString: The JWT token to be parsed.
//
// Returns the extracted username, expiration, and an error if extraction fails.
func ExtractClaimsFromToken(logger *logrus.Logger, tokenString string) (string, string, error) {
	mapClaims, err := extractMapClaimsFromToken(logger, tokenString)
	if err != nil {
		logger.WithError(err).WithField("token", tokenString).Error("Error extracting map claims from token")
		return "", "", err
	}

	username := fmt.Sprint(mapClaims["username"])
	exp := fmt.Sprint(mapClaims["exp"])
	logger.WithFields(logrus.Fields{
		"token":    tokenString,
		"username": username,
		"exp":      exp,
	}).Info("Token claims retrieved successfully")
	return username, exp, nil
}

// HandleTokensCreation generates and handles the creation of access and refresh tokens for a user.
// It creates a new access token and refresh token for the specified username and stores the refresh token in the database.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// db: A pointer to the SQL database connection.
// cfg: A pointer to the config.Config struct which contains JWT configuration.
// userName: The username for which tokens are being generated.
//
// Returns the generated access token, refresh token, and an error if token creation or storage fails.
func HandleTokensCreation(logger *logrus.Logger, db *sql.DB, cfg *config.Config, userName string) (string, string, error) {
	var accessToken, refreshToken string
	accessToken, err := createToken(logger, cfg, userName, cfg.JWTExpirationTime)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = createToken(logger, cfg, userName, cfg.JWTExpirationTime)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
