package database

import (
	"database/sql"
	"fmt"
)

// Creates a row in the Category table for this "category" structure.
func (category *Category) CreateCategory(db *sql.DB) (err error) {
	stmt, err := db.Prepare("INSERT INTO categories (name, post_id) values ($1, $2) returning id")
	err = stmt.QueryRow(category.Name, category.Post[0].Id).Scan(&category.Id)
	return
}

// Deletes a row in the Category table for this "category" structure.
func (category *Category) DeleteCategory(db *sql.DB) (err error) {
	_, err = db.Exec("DELETE FROM categories WHERE id = $1", category.Id)
	return
}

// Returns all categories of the post with the given ID.
func GetAllCategoryByPost(db *sql.DB, id int) ([]string, error) {
	statement := "SELECT name FROM categories WHERE post_id = $1"
	rows, err := db.Query(statement, id)
	if err != nil {
		return nil, fmt.Errorf("get all category by post: %w", err)
	}

	defer rows.Close()

	var result []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("get all category by post: %w", err)
		}
		result = append(result, category)
	}

	if rows.Err() != nil {
	}

	return result, nil
}

// Creates rows in the category table for given array of rows.
func CreateCategories(db *sql.DB, categories []string, post_id int) error {
	for _, category := range categories {
		_, err := db.Exec("INSERT INTO categories (name, post_id) VALUES ($1, $2)", category, post_id)
		if err != nil {
			return err
		}
	}

	return nil
}
