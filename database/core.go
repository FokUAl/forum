package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func Init() {
	var err error
	database, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		panic("failed to connect database")
	}

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS users " +
		"(id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, email TEXT)")
	statement.Exec()

	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS posts " +
		"(id INTEGER PRIMARY KEY, message TEXT, author TEXT, email TEXT, rating INTEGER, category_id Integer)")
	statement.Exec()

	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS users " +
		"(id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, email TEXT)")
	statement.Exec()
}
