package config

import (
	"os"
	"strconv"
)

func NewConfig() *Config {
	LoadEnv()

	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT")) // Convert port to int
	serverPort, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))

	return &Config{
		// Server Config
		ServerPort: serverPort,
		APIKey:     os.Getenv("API_KEY"),
		HashKey:    os.Getenv("HASH_KEY"),
		LogDir:     os.Getenv("LOG_DIR"),

		// DB Config
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     dbPort,
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}
}
