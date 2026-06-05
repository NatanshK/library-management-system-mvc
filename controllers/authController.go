package controllers

import (
	"encoding/json"
	"library-management-system-mvc/models"
	"library-management-system-mvc/utils"
	"net/http"
)

// Register handles new user sign-ups
func Register(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Parsing
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"error": "Invalid input data"}`, http.StatusBadRequest)
		return
	}

	// Generating a salt
	salt, err := utils.GenerateSalt(16)
	if err != nil {
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	user.Salt = salt
	user.Password = utils.HashPassword(user.Password, salt)

	// Default values for new users
	user.Role = "client"
	user.RequestStatus = "not_requested"

	// Saving to database
	err = models.CreateUser(&user)
	if err != nil {
		// If it fails here, it usually means the UNIQUE constraint on email/username caught a duplicate
		http.Error(w, `{"error": "Username or email already exists"}`, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully!"})
}
