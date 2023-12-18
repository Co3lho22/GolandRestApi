package config

type Config struct {
	//Server Config
	ServerPort int
	APIKey     string
	HashKey    string
	LogDir     string

	// DB Config
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}
