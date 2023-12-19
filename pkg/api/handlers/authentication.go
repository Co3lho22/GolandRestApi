package handlers

import (
	"GolandRestApi/pkg/model"
	"GolandRestApi/pkg/repository"
	"GolandRestApi/pkg/service"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RegisterUser(logger *logrus.Logger, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var newUser model.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		logger.WithError(err).Error("Failed to deserialize the User object for the /register endpoint")
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	userExists, err := repository.UserExists(logger, db, newUser.Username, newUser.Email)
	if err != nil {
		logger.WithError(err).Error("Error verifying if user exists")
		http.Error(w, "Error while verifying if user exists", http.StatusInternalServerError)
		return
	} else if userExists == true {
		http.Error(w, "Username or Email already in use", http.StatusMethodNotAllowed)
		return
	}

	hashedPassword, err := service.HashPassword(logger, newUser)
	if err != nil {
		logger.WithError(err).Error("Error adding user")
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	newUser.HashedPassword = hashedPassword
	success, err := repository.AddUser(logger, db, newUser)
	if err != nil {
		logger.WithError(err).Error("Error adding user")
		http.Error(w, "Error adding user", http.StatusInternalServerError)
		return
	}

	if !success {
		http.Error(w, "Error while registering user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("User successfully created"))
	if err != nil {
		logger.WithError(err).Error("Error writing the response")
		http.Error(w, "Error writing the response", http.StatusInternalServerError)
		return
	}

	logger.WithField("username", newUser.Username).Info("User registered with success")
	return
}

func LoginUser(logger *logrus.Logger, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	var loginDetails struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginDetails)
	if err != nil {
		logger.WithError(err).Error("Failed to deserialize the User object for the /login endpoint")
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	var newUser *model.User
	newUser, err = repository.GetUserByUserName(logger, db, loginDetails.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		logger.WithError(err).
			WithField("username", loginDetails.Username).
			Error("Error retrieving user from DB")
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	newUser.Username = loginDetails.Username
	if err := service.CheckPasswordHash(logger, newUser, loginDetails.Password); err != nil {
		logger.
			WithField("username", loginDetails.Username).
			Warn("Invalid password")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	responseMessage := "User successfully logged in"
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(responseMessage))
	if err != nil {
		logger.WithError(err).Error("Error writing the response")
		http.Error(w, "Error writing the response", http.StatusInternalServerError)
		return
	}

	logger.WithField("username", loginDetails.Username).Info("User logged in with success")
	return
}
