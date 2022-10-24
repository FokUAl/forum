package database

import "database/sql"

type Category struct {
	Id   int
	Name string
	Post []*Post
}

func (category *Category) CreateCategory(db *sql.DB) (err error) {
	stmt, err := db.Prepare("INSERT INTO categories (name, post_id) values ($1, $2) returning id")
	err = stmt.QueryRow(category.Name, category.Post[0].Id).Scan(&category.Id)
	return
}

func (category *Category) DeleteCategory(db *sql.DB) (err error) {
	_, err = db.Exec("DELETE FROM categories WHERE id = $1", category.Id)
	return
}
