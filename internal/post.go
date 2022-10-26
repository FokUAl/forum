package internal

import (
	"database/sql"
	"errors"
	"forumAA/database"
)

func CreatePost(db *sql.DB) error {
	// TO DO
	// Read data from front

	post := database.Post{}

	err := post.Create(db)
	return err
}

func EditPost(db *sql.DB) error {
	return errors.New("Some error")
}
