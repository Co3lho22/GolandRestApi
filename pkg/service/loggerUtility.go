package service

import (
	"GolandRestApi/pkg/config"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

func NewLogger(cfg *config.Config) (*logrus.Logger, error) {
	// Create log directory if it doesn't exist
	err := os.MkdirAll(cfg.LogDir, 0755)
	if err != nil {
		return nil, err
	}

	// Create log file
	logFilePath := filepath.Join(cfg.LogDir, "app.log")
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// Initialize logrus logger
	var logger = logrus.New()
	logger.Out = file
	logger.SetLevel(logrus.InfoLevel) // Set the log level
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Also log to stdout
	logger.SetOutput(io.MultiWriter(os.Stdout, file))

	return logger, nil
}
