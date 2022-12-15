package internal

import (
	"database/sql"
	"fmt"
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

// func ToLikePost(db *sql.DB, value int, post_id int) error {
// 	return nil
// }

func CountLikes(db *sql.DB, post_id int) (int, error) {
	likes, err := database.GetLikeByPost(db, post_id)
	if err != nil {
		return 0, fmt.Errorf("CountLikes: %w", err)
	}

	result := 0
	for _, like := range likes {
		result += like.Value
	}

	return result, nil
}
