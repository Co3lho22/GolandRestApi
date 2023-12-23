package service

import (
	"GolandRestApi/pkg/utils"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

// HttpErrorResponse sends an HTTP error response to the client and logs the error message and details.
//
// logger: A logrus.Logger instance for logging information, warnings, and errors.
// w: The http.ResponseWriter to write the error response to.
// statusCode: The HTTP status code to set in the response.
// endpoint: The endpoint or URL path where the error occurred.
// message: The error message to include in the response.
// err: The error object or nil if there's no specific error.
// logType: The type of log message (e.g., utils.LogTypeInfo, utils.LogTypeWarn, utils.LogTypeError).
// username: The username associated with the request (used for logging purposes).
//
// This function logs the error message and details based on the logType and sends an HTTP response
// with the specified status code and error message in JSON format to the client.
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
