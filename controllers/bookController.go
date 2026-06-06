package controllers

import (
	"encoding/json"
	"library-management-system-mvc/models"
	"library-management-system-mvc/utils"

	//"log"
	"net/http"
	"strconv"
)

func AddBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	book.AvailableCopies = book.TotalCopies

	err = models.AddBook(&book)
	if err != nil {
		//log.Println("BOOK DATABASE REJECTION REASON:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add book to catalog")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Book added successfully!"})
}

func GetCatalog(w http.ResponseWriter, r *http.Request) {

	searchTerm := r.URL.Query().Get("q")

	books, err := models.GetBooks(searchTerm)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve catalog")
		return
	}

	if books == nil {
		books = []models.Book{}
	}

	utils.RespondWithJSON(w, http.StatusOK, books)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {

	idParam := r.URL.Query().Get("id")
	bookID, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Valid Book ID is required in the URL")
		return
	}

	var book models.Book
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	book.ID = bookID

	err = models.UpdateBook(&book)
	if err != nil {
		//log.Println("UPDATE DATABASE REJECTION REASON:", err)
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Book updated successfully!"})
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {

	idParam := r.URL.Query().Get("id")
	bookID, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Valid Book ID is required in the URL")
		return
	}

	err = models.DeleteBook(bookID)
	if err != nil {

		utils.RespondWithError(w, http.StatusConflict, "Cannot delete book: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Book deleted successfully!"})
}
