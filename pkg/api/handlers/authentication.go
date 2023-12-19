package handlers

import (
	"GolandRestApi/pkg/model"
	"GolandRestApi/pkg/repository"
	"GolandRestApi/pkg/service"
	"database/sql"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RegisterUser(logger *logrus.Logger, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var newUser model.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		logger.Errorf("Was not able to deserialize the User object for the /register endpoint")
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	//TODO: verify if a user with the same user name or email exists if yes respond to the client
	userExists, err := repository.UserExists(logger, db, newUser.Username, newUser.Email)
	if err != nil {
		logger.Errorf("Verifying if user exists")
		http.Error(w, "Verifying if user exists", http.StatusInternalServerError)
		return
	} else if userExists == true {
		http.Error(w, "Username or Email already in use", http.StatusMethodNotAllowed)
		return
	}

	hashedPassword, err := service.HashPassword(logger, newUser)
	if err != nil {
		logger.Errorf("Error hashing password")
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	newUser.HashedPassword = hashedPassword

	//TODO: Register user in the DB
	success, err := repository.AddUser(logger, db, newUser)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	if success {
		_, err := w.Write([]byte("User successfully created"))
		if err != nil {
			http.Error(w, "Error writing the response", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Error while registering user", http.StatusInternalServerError)
		return
	}

}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("This is the login endpoint"))
	if err != nil {
		return
	}
}
