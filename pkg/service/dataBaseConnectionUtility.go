package service

import (
	"GolandRestApi/pkg/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func NewDBConnection(logger *logrus.Logger, cfg *config.Config) (*sql.DB, error) {
	// Construct the connection string
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Open a new connection to the database
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	logger.Info("Connected to the database successfully")
	return db, nil
}
