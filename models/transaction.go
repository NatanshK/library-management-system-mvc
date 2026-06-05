package models

import "time"

type Transaction struct {
	TransactionID int        `json:"transaction_id"`
	UserID        int        `json:"user_id"`
	BookID        int        `json:"book_id"`
	Status        string     `json:"status"`
	CheckoutTime  *time.Time `json:"checkout_time"`
	CheckinTime   *time.Time `json:"checkin_time"`
	DueDate       *time.Time `json:"due_date"`
	FineAmount    float64    `json:"fine_amount"`
}
