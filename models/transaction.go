package models

import (
	"fmt"
	"library-management-system-mvc/config"
	"library-management-system-mvc/utils"
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
	BookID        int    `json:"book_id"`
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
		SET status = 'checkout_accepted', checkout_time = ?, due_date = ? 
		WHERE transaction_id = ? AND status = 'checkout_requested'
	`
	txResult, err := tx.Exec(updateTxQuery, now, dueDate, transactionID)
	if err != nil {
		tx.Rollback()
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
			t.transaction_id, b.title, b.author, t.status, 
			t.checkout_time, t.checkin_time, t.due_date, t.fine_amount
		FROM transactions t
		INNER JOIN books b ON t.book_id = b.id
		WHERE t.user_id = ?
		ORDER BY t.transaction_id DESC
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
			t.transaction_id, u.username, u.email, t.book_id, b.title, t.status
		FROM transactions t
		INNER JOIN users u ON t.user_id = u.id
		INNER JOIN books b ON t.book_id = b.id
		WHERE t.status IN ('checkout_requested', 'checkin_requested')
		ORDER BY t.transaction_id ASC
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
			&req.BookID,
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

func RequestCheckin(transactionID int, userID int) error {
	query := `
		UPDATE transactions 
		SET status = 'checkin_requested' 
		WHERE transaction_id = ? AND user_id = ? AND status = 'checkout_accepted'
	`

	res, err := config.DB.Exec(query, transactionID, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("invalid request: no active checkout found for this transaction")
	}

	return nil
}

func ApproveCheckin(transactionID int) (float64, error) {

	tx, err := config.DB.Begin()
	if err != nil {
		return 0.0, err
	}

	var dueDate *time.Time
	var bookID int
	readQuery := `
		SELECT due_date, book_id 
		FROM transactions 
		WHERE transaction_id = ? AND status = 'checkin_requested'
	`
	err = tx.QueryRow(readQuery, transactionID).Scan(&dueDate, &bookID)
	if err != nil {
		tx.Rollback()
		return 0.0, fmt.Errorf("transaction not found or not in checkin_requested state")
	}

	now := time.Now()
	var fineAmount float64 = 0.0

	if dueDate != nil {
		_, fineAmount = utils.CalculateFine(*dueDate, now, 2.00)
	}

	updateTxQuery := `
		UPDATE transactions 
		SET status = 'returned', checkin_time = ?, fine_amount = ? 
		WHERE transaction_id = ?
	`
	_, err = tx.Exec(updateTxQuery, now, fineAmount, transactionID)
	if err != nil {
		tx.Rollback()
		return 0.0, err
	}

	updateBookQuery := `
		UPDATE books 
		SET available_copies = available_copies + 1 
		WHERE id = ?
	`
	_, err = tx.Exec(updateBookQuery, bookID)
	if err != nil {
		tx.Rollback()
		return 0.0, err
	}

	return fineAmount, tx.Commit()
}

func RejectCheckout(transactionID int) error {
	query := `
		UPDATE transactions 
		SET status = 'checkout_rejected' 
		WHERE transaction_id = ? AND status = 'checkout_requested'
	`

	res, err := config.DB.Exec(query, transactionID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("cannot reject: transaction not found or not in checkout_requested state")
	}

	return nil
}

func RejectCheckin(transactionID int) error {

	query := `
		UPDATE transactions 
		SET status = 'checkout_accepted' 
		WHERE transaction_id = ? AND status = 'checkin_requested'
	`

	res, err := config.DB.Exec(query, transactionID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("cannot reject: transaction not found or not in checkin_requested state")
	}

	return nil
}
