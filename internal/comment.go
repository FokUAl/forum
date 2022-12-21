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

func CountCommentLikes(db *sql.DB, comment_id int) (int, error) {
	result := 0

	likes, err := database.GetLikeByComment(db, comment_id)
	if err != nil {
		return result, err
	}

	for _, like := range likes {
		result += like.Value
	}

	return result, nil
}
