package config

// Config represents the configuration settings for the GoLandRestApi application.
type Config struct {
	// Server Configuration
	ServerPort int

	// Database Configuration
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// Logging Configuration
	LogDir string

	// JWT Configuration
	JWTSecretKey            string
	JWTExpirationTime       string
	JWTRefreshTokenValidity string
}
