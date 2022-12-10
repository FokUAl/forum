package database

import (
	"database/sql"
	"fmt"
)

func (category *Category) CreateCategory(db *sql.DB) (err error) {
	stmt, err := db.Prepare("INSERT INTO categories (name, post_id) values ($1, $2) returning id")
	err = stmt.QueryRow(category.Name, category.Post[0].Id).Scan(&category.Id)
	return
}

func (category *Category) DeleteCategory(db *sql.DB) (err error) {
	_, err = db.Exec("DELETE FROM categories WHERE id = $1", category.Id)
	return
}

func GetAllCategoryByPost(db *sql.DB, id int) ([]string, error) {
	statement := "SELECT name FROM categories WHERE post_id = $1"
	rows, err := db.Query(statement, id)
	if err != nil {
		return nil, fmt.Errorf("get all category by post: %w", err)
	}

	var result []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("get all category by post: %w", err)
		}
		result = append(result, category)
	}

	return result, nil
}
