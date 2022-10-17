package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func Init() {
	var err error
	database, err = sql.Open("sqlite3", "./forum.db?_foreign_keys=on")
	if err != nil {
		panic("failed to connect database")
	}

	/*
		statement, _ := database.Prepare("PRAGMA foreign_keys = on")
		statement.Exec()
	*/

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS users " +
		"(id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, nickname TEXT," +
		"email TEXT, post_id INTEGER)" +
		"FOREIGN KEY (post_id) REFERENCES posts(id)" +
		"CONSTRAINT name_unique UNIQUE(nickname)")
	statement.Exec()

	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS posts " +
		"(id INTEGER PRIMARY KEY, message TEXT, author TEXT, email TEXT, " +
		"like INTEGER, dislike INTEGER, category_id INTEGER)" +
		"FOREIGN KEY (category_id) REFERENCES categories(id)" +
		"FOREIGN KEY (author) REFERENCES users(nickname)" +
		"ON DELETE CASCADE")
	statement.Exec()

	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS categories " +
		"(id INTEGER PRIMARY KEY, name TEXT, post_id INTEGER,)" +
		"FOREIGN KEY (post_id) REFERENCES posts(id)" +
		"ON DELETE CASCADE")
	statement.Exec()
}
