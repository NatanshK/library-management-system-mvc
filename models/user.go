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

func RequestAdminRole(userID int) error {

	query := `
		UPDATE users 
		SET request_status = 'pending' 
		WHERE id = ? AND role != 'admin' AND request_status != 'pending'
	`

	res, err := config.DB.Exec(query, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("cannot apply: user not found, already an admin, or request already pending")
	}

	return nil
}

func ApproveAdminRole(userID int) error {
	query := `
		UPDATE users 
		SET role = 'admin', request_status = 'accepted' 
		WHERE id = ? AND request_status = 'pending'
	`

	res, err := config.DB.Exec(query, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("cannot approve: user not found or no pending request exists")
	}

	return nil
}

func GetPendingAdmins() ([]User, error) {
	query := `SELECT id, username, email, role, request_status FROM users WHERE request_status = 'pending'`

	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User

		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.RequestStatus)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
