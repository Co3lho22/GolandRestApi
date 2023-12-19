package service

import (
	"GolandRestApi/pkg/config"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

// NewLogger creates a new logrus.Logger instance for logging application messages and errors.
// It takes a pointer to the config.Config struct, which contains the configuration details
// necessary for setting up the logger.
//
// cfg: A pointer to the config.Config struct which specifies the log directory path and other
// configuration parameters.
//
// Returns a pointer to a logrus.Logger object and an error. If there is an error in creating
// or configuring the logger, it returns nil and the error.
func NewLogger(cfg *config.Config) (*logrus.Logger, error) {
	err := os.MkdirAll(cfg.LogDir, 0755)
	if err != nil {
		return nil, err
	}

	logFilePath := filepath.Join(cfg.LogDir, "app.log")
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	var logger = logrus.New()
	logger.Out = file
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logger.SetOutput(io.MultiWriter(os.Stdout, file))

	return logger, nil
}
