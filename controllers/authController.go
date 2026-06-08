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

type RoleRequest struct {
	UserID int `json:"user_id"`
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

func RequestPromotion(w http.ResponseWriter, r *http.Request) {
	var req RoleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := models.RequestAdminRole(req.UserID); err != nil {
		utils.RespondWithError(w, http.StatusConflict, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Admin promotion requested successfully. Pending review.",
	})
}

func ApprovePromotion(w http.ResponseWriter, r *http.Request) {
	var req RoleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := models.ApproveAdminRole(req.UserID); err != nil {
		utils.RespondWithError(w, http.StatusConflict, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Promotion approved. User is now an admin.",
	})
}

func GetPendingPromotions(w http.ResponseWriter, r *http.Request) {
	queue, err := models.GetPendingAdmins()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve HR queue")
		return
	}

	if queue == nil {
		queue = []models.User{}
	}

	utils.RespondWithJSON(w, http.StatusOK, queue)
}
