package internal

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
)

func Init(){
	database, _ := sql.Open(driverName: "sqlite3", dataSourceName:"./forum.db")
	statement, _ := database.Prepare(query: "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, email TEXT)")
	statement.Exec()
}


