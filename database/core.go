package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Init() *sql.DB {
	var err error
	db, err := sql.Open("sqlite3", "file:forum.db")
	if err != nil {
		panic("failed to connect database")
	}

	statement, err := db.Prepare("PRAGMA foreign_keys = 1")
	if err != nil {
		log.Fatal("Init: " + err.Error())
	}
	statement.Exec()

	// CreateTables()
	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS users " +
		"(id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, nickname TEXT UNIQUE," +
		"password TEXT, email TEXT UNIQUE)")
	if err != nil {
		log.Fatal("Init: users: " + err.Error())
	}
	statement.Exec()

	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS posts " +
		"(id INTEGER PRIMARY KEY, title TEXT, message TEXT, author TEXT, " +
		"user_id INTEGER, FOREIGN KEY (user_id) REFERENCES users(id))")
	if err != nil {
		log.Fatal("Init: posts: " + err.Error())
	}
	statement.Exec()

	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS categories " +
		"(id INTEGER PRIMARY KEY, name TEXT, post_id INTEGER, " +
		"FOREIGN KEY (post_id) REFERENCES posts(id) " +
		"ON DELETE CASCADE)")
	if err != nil {
		log.Fatalf("Init: categories: " + err.Error())
	}
	statement.Exec()

	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS comments " +
		"(id INTEGER PRIMARY KEY, content TEXT, author TEXT, post_id INTEGER," +
		"FOREIGN KEY (post_id) REFERENCES posts(id) " +
		"ON DELETE CASCADE)")
	if err != nil {
		log.Fatal("Init: comments: " + err.Error())
	}
	statement.Exec()

	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS sessions " +
		"(id INTEGER PRIMARY KEY, nickname TEXT UNIQUE, session_token TEXT, expiry DATETIME)")
	if err != nil {
		log.Fatalf("Init: sessions: " + err.Error())
	}
	statement.Exec()

	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS comment_likes " +
		"(id INTEGER PRIMARY KEY, nickname TEXT UNIQUE, like INTEGER DEFAULT 0, comment_id INTEGER," +
		"FOREIGN KEY (comment_id) REFERENCES comments(id) " +
		"ON DELETE CASCADE)")
	if err != nil {
		log.Fatalf("Init: comments_likes: " + err.Error())
	}
	statement.Exec()

	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS post_likes " +
		"(id INTEGER PRIMARY KEY, nickname TEXT UNIQUE, like INTEGER DEFAULT 0, post_id INTEGER, " +
		"FOREIGN KEY (post_id) REFERENCES posts(id) " +
		"ON DELETE CASCADE)")
	if err != nil {
		log.Fatalf("Init: post_likes: " + err.Error())
	}
	statement.Exec()

	return db
}
