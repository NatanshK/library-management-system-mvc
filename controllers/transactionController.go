package controllers

import (
	"encoding/json"
	"library-management-system-mvc/models"
	"library-management-system-mvc/utils"

	//"log"
	"net/http"
	"strconv"
)

type TransactionRequest struct {
	UserID        int `json:"user_id,omitempty"`
	BookID        int `json:"book_id,omitempty"`
	TransactionID int `json:"transaction_id,omitempty"`
}

func RequestCheckout(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	err = models.RequestCheckout(req.UserID, req.BookID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to submit checkout request")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "Checkout requested successfully. Pending admin approval.",
	})
}

func ApproveCheckout(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	err = models.ApproveCheckout(req.TransactionID, req.BookID)
	if err != nil {
		utils.RespondWithError(w, http.StatusConflict, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Checkout approved. Book inventory updated.",
	})
}

func GetUserHistory(w http.ResponseWriter, r *http.Request) {

	userIDParam := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Valid User ID is required")
		return
	}

	history, err := models.GetUserHistory(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user history")
		return
	}

	if history == nil {
		history = []models.UserHistoryDTO{}
	}

	utils.RespondWithJSON(w, http.StatusOK, history)
}

func GetAdminQueue(w http.ResponseWriter, r *http.Request) {
	queue, err := models.GetPendingRequests()
	if err != nil {
		// log.Println("ADMIN QUEUE DATABASE REJECTION REASON:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve admin queue")
		return
	}

	if queue == nil {
		queue = []models.AdminQueueDTO{}
	}

	utils.RespondWithJSON(w, http.StatusOK, queue)
}
