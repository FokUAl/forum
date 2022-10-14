package internal

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	Id         int
	Message    string
	Author     string
	rating     int
	Categories []string
}

var database *sql.DB

func Init() {
	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		panic("failed to connect database")
	}
}

func (post *Post) Create() (err error) {
	statement := "INSERT INTO posts (message, author) values ($1, $2) returning id"
	stmt, err := database.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(post.Message, post.Author).Scan(&post.Id)
	return
}
