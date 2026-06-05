package models

import (
	"fmt"
	"library-management-system-mvc/config"
)

type Book struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	ISBN            string `json:"isbn"`
	PublicationYear int    `json:"publication_year"`
	TotalCopies     int    `json:"total_copies"`
	AvailableCopies int    `json:"available_copies"`
}

func AddBook(book *Book) error {
	query := `
		INSERT INTO books (title, author, isbn, publication_year, total_copies, available_copies) 
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := config.DB.Exec(query,
		book.Title,
		book.Author,
		book.ISBN,
		book.PublicationYear,
		book.TotalCopies,
		book.AvailableCopies,
	)

	return err
}

func GetBooks(searchTerm string) ([]Book, error) {
	var books []Book

	query := `
		SELECT id, title, author, isbn, publication_year, total_copies, available_copies 
		FROM books 
	`
	var args []interface{}

	if searchTerm != "" {

		query += " WHERE title LIKE ? OR author LIKE ?"

		wildcardSearch := "%" + searchTerm + "%"

		args = append(args, wildcardSearch, wildcardSearch)
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	//Ensure the connection is released back to the pool when done
	defer rows.Close()

	for rows.Next() {
		var b Book

		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Author,
			&b.ISBN,
			&b.PublicationYear,
			&b.TotalCopies,
			&b.AvailableCopies,
		)
		if err != nil {
			return nil, err
		}

		books = append(books, b)
	}

	return books, nil
}

func UpdateBook(book *Book) error {
	query := `
		UPDATE books 
		SET title = ?, author = ?, isbn = ?, publication_year = ?, total_copies = ?, available_copies = ?
		WHERE id = ?
	`

	result, err := config.DB.Exec(query,
		book.Title,
		book.Author,
		book.ISBN,
		book.PublicationYear,
		book.TotalCopies,
		book.AvailableCopies,
		book.ID,
	)
	if err != nil {
		return err
	}

	//The "Ghost Update"
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no book found with ID %d", book.ID)
	}

	return nil
}

// MySQL's foreign key constraints prevents deleting actively borrowed books.
func DeleteBook(id int) error {
	query := "DELETE FROM books WHERE id = ?"

	result, err := config.DB.Exec(query, id)

	if err != nil {
		// If MySQL blocks the deletion because a student is borrowing it
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cannot delete: no book found with ID %d", id)
	}

	return nil
}
