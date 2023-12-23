package service

import (
	"GolandRestApi/pkg/utils"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func HttpErrorResponse(logger *logrus.Logger,
	w http.ResponseWriter,
	statusCode int,
	endpoint,
	message string,
	err error,
	logType,
	username string) {

	switch logType {
	case utils.LogTypeInfo:
		logger.WithFields(logrus.Fields{
			"endpoint": endpoint,
			"username": username,
		}).Info(message)
	case utils.LogTypeWarn:
		logger.WithFields(logrus.Fields{
			"endpoint": endpoint,
			"username": username,
		}).Warn(message)
	case utils.LogTypeError:
		logger.WithError(err).WithFields(logrus.Fields{
			"endpoint": endpoint,
			"username": username,
		}).Error(message)
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
