package service

import (
	"GolandRestApi/pkg/config"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"time"
)

type Claims struct {
	UserId int `json:"userId"`
	jwt.StandardClaims
}

// TODO: Add method to check the refresh token

// TODO: Add method to generate token

func GenerateAccessToken(logger *logrus.Logger,
	cfg *config.Config,
	userId int,
	username string,
	roles []string,
) (string, error) {

	// Parse duration from config
	expirationTime, err := time.ParseDuration(cfg.JWTExpirationTime)
	if err != nil {
		logger.WithError(err).Error("Invalid JWT expiration time format")
		return "", err
	}

	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			Subject:   username,
			ExpiresAt: time.Now().Add(expirationTime).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecretKey))
	if err != nil {
		logger.Errorf("Failed to generate JWT token: %v", err)
		return "", err
	}

	logger.Infof("JWT token generated for user: %s", username)
	return tokenString, nil
}

func GenerateRefreshToken(logger *logrus.Logger, cfg *config.Config, userId int, username string, roles []string) (string, error) {
	refreshTokenValidity, err := time.ParseDuration(cfg.JWTRefreshTokenValidity)
	if err != nil {
		logger.WithError(err).Error("Invalid JWT refresh token validity format")
		return "", err
	}

	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			Subject:   username,
			ExpiresAt: time.Now().Add(refreshTokenValidity).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	refreshToken, err := token.SignedString([]byte(cfg.JWTSecretKey))
	if err != nil {
		logger.Errorf("Failed to generate refresh token: %v", err)
		return "", err
	}

	logger.Infof("Refresh token generated for user: %s", username)
	return refreshToken, nil
}

// TODO: Add method to verify token
func VerifyToken(logger *logrus.Logger, cfg *config.Config, tokenString string) (int, string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecretKey), nil
	})

	if err != nil {
		logger.WithError(err).Error("Error verifying JWT token")
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	logger.Infof("JWT token verified successfully for user: %s", claims.Subject)
	return claims.UserId, claims.Subject, nil
}

func getUsernameFromToken(logger *logrus.Logger, cfg *config.Config, tokenString string) (string, error) {
	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecretKey), nil
	})

	if err != nil {
		logger.WithError(err).Error("Error verifying JWT token")
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	logger.Infof("JWT token verified successfully for user: %s", claims.Subject)
	return claims.Subject, nil
}

func getUserIdFromToken(logger *logrus.Logger, cfg *config.Config, tokenString string) (int, error) {
	claims := &jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecretKey), nil
	})

	if err != nil {
		logger.WithError(err).Error("Error verifying JWT token")
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	userId, ok := (*claims)["userId"].(float64) // JWT numeric values are float64
	if !ok {
		return 0, errors.New("user ID not found in token")
	}

	logger.Infof("JWT token verified successfully for user ID: %d", int(userId))
	return int(userId), nil
}
