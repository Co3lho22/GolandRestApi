package config

import (
	"os"
	"strconv"
)

// NewConfig creates a new Config instance by loading configuration values from environment variables.
// It reads environment variables for various configuration parameters, converts them to their respective types,
// and initializes the Config struct accordingly.
//
// The following environment variables are used:
// - SERVER_PORT: The port on which the HTTP server should listen.
// - DB_HOST: The hostname or IP address of the MySQL database server.
// - DB_PORT: The port on which the MySQL database server is running.
// - DB_USER: The username used to authenticate with the MySQL database server.
// - DB_PASSWORD: The password used to authenticate with the MySQL database server.
// - DB_NAME: The name of the MySQL database to connect to.
// - LOG_DIR: The directory where log files should be stored.
// - JWT_SECRET_KEY: The secret key used for JWT token signing and verification.
// - JWT_EXPIRATION_TIME: The duration of JWT token validity (e.g., "15m" for 15 minutes).
// - JWT_REFRESH_TOKEN_VALIDITY: The duration of JWT refresh token validity (e.g., "7d" for 7 days).
//
// Returns a pointer to a Config struct initialized with the loaded configuration values.
func NewConfig() *Config {
	LoadEnv()

	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT")) // Convert port to int
	serverPort, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))

	return &Config{
		// Server Configuration
		ServerPort: serverPort,

		// Database Configuration
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     dbPort,
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		// Logging Configuration
		LogDir: os.Getenv("LOG_DIR"),

		// JWT Configuration
		JWTSecretKey:            os.Getenv("JWT_SECRET_KEY"),
		JWTExpirationTime:       os.Getenv("JWT_EXPIRATION_TIME"),
		JWTRefreshTokenValidity: os.Getenv("JWT_REFRESH_TOKEN_VALIDITY"),
	}
}
