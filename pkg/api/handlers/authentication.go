package handlers

import (
	"GolandRestApi/pkg/model"
	"encoding/json"
	"net/http"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser model.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		return
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("This is the login endpoint"))
	if err != nil {
		return
	}
}
