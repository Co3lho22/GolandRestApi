package config

import (
	"log"
	"os"
	"strconv"
)

// NewConfig creates a new configuration instance for the application based on environment variables.
// It reads various environment variables to configure the application, such as database settings,
// server settings, and JWT token settings. If any of the required environment variables are not set
// or have invalid values, it will log a fatal error and terminate the application.
//
// Returns a pointer to a Config struct containing the configuration parameters for the application.
func NewConfig() *Config {
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("DB_PORT environment variable is not set or invalid")
	}

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatalf("SERVER_PORT environment variable is not set or invalid")
	}

	return &Config{
		ServerPort:              serverPort,
		APIVersion:              getEnv("API_VERSION", "v1"),
		DBHost:                  getEnv("DB_HOST", "localhost"),
		DBPort:                  dbPort,
		DBUser:                  getEnv("DB_USER", "defaultUser"),
		DBPassword:              getEnv("DB_PASSWORD", ""),
		DBName:                  getEnv("DB_NAME", "RestApi"),
		LogDir:                  getEnv("LOG_DIR", "/var/log/restapi/"),
		JWTSecretKey:            getEnv("JWT_SECRET_KEY", "defaultSecret"),
		JWTExpirationTime:       getEnv("JWT_EXPIRATION_TIME", "15m"),
		JWTRefreshTokenValidity: getEnv("JWT_REFRESH_TOKEN_VALIDITY", "7d"),
	}
}

// getEnv retrieves the value of an environment variable specified by the 'key' parameter.
// If the environment variable is not set, it returns the 'fallback' value.
//
// key: The name of the environment variable to retrieve.
// fallback: The default value to return if the environment variable is not set.
//
// Returns a string representing the value of the environment variable or the fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
