package internal

import (
	"database/sql"
	"forumAA/database"
)

func CreateComment(db *sql.DB) error {
	// TO DO
	// Read data from front
	comment := database.Comment{}

	return comment.Create(db)
}
