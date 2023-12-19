package service

import (
	"GolandRestApi/pkg/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// NewDBConnection establishes a new connection to the MySQL database using the provided configuration.
// It constructs the connection string using the database user, password, host, port, and name from the config.
// After opening a new database connection, it attempts to ping the database to verify the connection.
// On success, it returns a pointer to the sql.DB instance and nil error.
// If there is any error in opening the connection or during the ping, the function returns nil and the error.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// cfg: A pointer to the config.Config struct which contains the database configuration details.
//
// Returns a pointer to a sql.DB object and an error.
func NewDBConnection(logger *logrus.Logger, cfg *config.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	logger.Info("Connected to the database successfully")
	return db, nil
}
