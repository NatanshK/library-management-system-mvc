package controllers

import (
	"encoding/json"
	"library-management-system-mvc/models"
	"library-management-system-mvc/utils"

	//"log"
	"net/http"
)

type Credentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	user := models.User{
		Username:      creds.Username,
		Email:         creds.Email,
		Password:      creds.Password,
		Role:          "client",
		RequestStatus: "not_requested",
	}

	err = models.CreateUser(&user)
	if err != nil {
		//log.Println("DATABASE REJECTION REASON:", err)
		utils.RespondWithError(w, http.StatusConflict, "Email or Username already exists")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully!"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	user, err := models.GetUserByEmail(creds.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	hashedInput := utils.HashPassword(creds.Password, user.Salt)

	if hashedInput != user.Password {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate session token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}
