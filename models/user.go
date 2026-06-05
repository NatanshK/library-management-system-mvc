package models

import (
	"database/sql"
	"errors"
	"library-management-system-mvc/config"
	"library-management-system-mvc/utils"
)

type User struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	Password      string `json:"-"` // Hidden from API responses
	Email         string `json:"email"`
	Role          string `json:"role"`
	RequestStatus string `json:"request_status"`
	Salt          string `json:"-"` // Hidden from API responses
}

func CreateUser(user *User) error {
	salt, err := utils.GenerateSalt(16)
	if err != nil {
		return err
	}

	hashedPassword := utils.HashPassword(user.Password, salt)

	query := `
		INSERT INTO users (username, password, email, role, salt) 
		VALUES (?, ?, ?, ?, ?)
	`

	_, err = config.DB.Exec(query, user.Username, hashedPassword, user.Email, user.Role, salt)
	if err != nil {
		return err
	}

	return nil
}
func GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, username, password, email, role, request_status, salt FROM users WHERE email = ?`

	var user User

	err := config.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Role,
		&user.RequestStatus,
		&user.Salt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
