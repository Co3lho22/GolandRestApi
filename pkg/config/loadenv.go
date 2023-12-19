package config

import (
	"github.com/joho/godotenv"
	"log"
)

// LoadEnv loads environment variables from a .env file located in the application's directory.
// It uses the "github.com/joho/godotenv" package to read the .env file and set environment variables.
//
// If the .env file is not found or there is an error while loading it, the function logs an error and terminates
// the application.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
