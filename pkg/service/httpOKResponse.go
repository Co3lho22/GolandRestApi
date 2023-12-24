package service

import (
	"GolandRestApi/pkg/utils"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func HttpOKResponse(logger *logrus.Logger,
	w http.ResponseWriter,
	endpoint,
	message string,
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
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": message})

}
