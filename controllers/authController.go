package controllers

import (
	"encoding/json"
	"library-management-system-mvc/models"
	"library-management-system-mvc/utils"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"error": "Invalid input data"}`, http.StatusBadRequest)
		return
	}

	salt, err := utils.GenerateSalt(16)
	if err != nil {
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	user.Salt = salt
	user.Password = utils.HashPassword(user.Password, salt)

	user.Role = "client"
	user.RequestStatus = "not_requested"

	err = models.CreateUser(&user)
	if err != nil {
		http.Error(w, `{"error": "Username or email already exists"}`, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully!"})
}
