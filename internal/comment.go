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

func EditComment(db *sql.DB) error {
	// TO DO
	// READ data from front
	id := 1
	newContent := ""

	comment, err := database.GetComment(db, id)
	if err != nil {
		return err
	}

	err = comment.Update(db, newContent)
	return err
}
