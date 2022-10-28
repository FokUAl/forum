package internal

import (
	"database/sql"
	"forumAA/database"
)

func CreatePost(db *sql.DB) error {
	// TO DO
	// Read data from front

	post := database.Post{}

	err := post.Create(db)
	return err
}

func EditPost(db *sql.DB, title, message string) error {
	// TO DO
	// Read id from post
	id := 1
	post, err := database.GetPost(db, id)
	if err != nil {
		return err
	}

	err = post.Update(db, title, message)
	return err
}
