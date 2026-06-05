package models

import (
	"fmt"
	"library-management-system-mvc/config"
	"time"
)

type Transaction struct {
	ID           int        `json:"transaction_id"`
	UserID       int        `json:"user_id"`
	BookID       int        `json:"book_id"`
	Status       string     `json:"status"`
	CheckoutTime *time.Time `json:"checkout_time"`
	CheckinTime  *time.Time `json:"checkin_time"`
	DueDate      *time.Time `json:"due_date"`
	FineAmount   float64    `json:"fine_amount"`
}

type UserHistoryDTO struct {
	TransactionID int        `json:"transaction_id"`
	BookTitle     string     `json:"book_title"`
	BookAuthor    string     `json:"book_author"`
	Status        string     `json:"status"`
	CheckoutTime  *time.Time `json:"checkout_time"`
	CheckinTime   *time.Time `json:"checkin_time"`
	DueDate       *time.Time `json:"due_date"`
	FineAmount    float64    `json:"fine_amount"`
}

type AdminQueueDTO struct {
	TransactionID int    `json:"transaction_id"`
	UserName      string `json:"user_name"`
	UserEmail     string `json:"user_email"`
	BookTitle     string `json:"book_title"`
	Status        string `json:"status"`
}

func RequestCheckout(userID int, bookID int) error {
	query := `INSERT INTO transactions (user_id, book_id, status) VALUES (?, ?, 'checkout_requested')`
	_, err := config.DB.Exec(query, userID, bookID)
	return err
}

// ACID Transaction
func ApproveCheckout(transactionID int, bookID int) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}

	now := time.Now()
	dueDate := now.AddDate(0, 0, 14)

	updateTxQuery := `
		UPDATE transactions 
		SET status = 'borrowed', checkout_time = ?, due_date = ? 
		WHERE id = ? AND status = 'checkout_requested'
	`
	txResult, err := tx.Exec(updateTxQuery, now, dueDate, transactionID)
	if err != nil {
		tx.Rollback() // rollback if it crashes in between
		return err
	}

	txRows, _ := txResult.RowsAffected()
	if txRows == 0 {
		tx.Rollback()
		return fmt.Errorf("transaction %d not found or already processed", transactionID)
	}

	updateBookQuery := `
		UPDATE books 
		SET available_copies = available_copies - 1 
		WHERE id = ? AND available_copies > 0
	`
	bookResult, err := tx.Exec(updateBookQuery, bookID)
	if err != nil {
		tx.Rollback()
		return err
	}

	bookRows, _ := bookResult.RowsAffected()
	if bookRows == 0 {
		tx.Rollback()
		return fmt.Errorf("cannot approve: book %d is currently out of stock", bookID)
	}

	return tx.Commit()
}

func GetUserHistory(userID int) ([]UserHistoryDTO, error) {
	var history []UserHistoryDTO

	query := `
		SELECT 
			t.id, b.title, b.author, t.status, 
			t.checkout_time, t.checkin_time, t.due_date, t.fine_amount
		FROM transactions t
		INNER JOIN books b ON t.book_id = b.id
		WHERE t.user_id = ?
		ORDER BY t.id DESC
	`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var h UserHistoryDTO

		err := rows.Scan(
			&h.TransactionID,
			&h.BookTitle,
			&h.BookAuthor,
			&h.Status,
			&h.CheckoutTime,
			&h.CheckinTime,
			&h.DueDate,
			&h.FineAmount,
		)
		if err != nil {
			return nil, err
		}

		history = append(history, h)
	}

	return history, nil
}

func GetPendingRequests() ([]AdminQueueDTO, error) {
	var queue []AdminQueueDTO

	query := `
		SELECT 
			t.id, u.username, u.email, b.title, t.status
		FROM transactions t
		INNER JOIN users u ON t.user_id = u.id
		INNER JOIN books b ON t.book_id = b.id
		WHERE t.status IN ('checkout_requested', 'return_requested')
		ORDER BY t.id ASC
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var req AdminQueueDTO

		err := rows.Scan(
			&req.TransactionID,
			&req.UserName,
			&req.UserEmail,
			&req.BookTitle,
			&req.Status,
		)
		if err != nil {
			return nil, err
		}

		queue = append(queue, req)
	}

	return queue, nil
}
