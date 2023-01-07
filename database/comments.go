package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// Creates a row in the Comments table for this "comment" structure.
func (comment *Comment) Create(db *sql.DB) (err error) {
	if comment.Post_Id == 0 {
		err = errors.New("Post not found")
		return
	}
	statement := "INSERT INTO comments (content, author, post_id)" +
		"values ($1, $2, $3) returning id"
	err = db.QueryRow(statement, comment.Content, comment.Author, comment.Post_Id).Scan(&comment.Id)
	return
}

// Deletes a row in the Comments table for this "comment" structure.
func (comment *Comment) Delete(db *sql.DB) (err error) {
	_, err = db.Exec("DELETE FROM comments where id = $1", comment.Id)
	return
}

// Returns the "comment" structure from the database by its ID.
func GetComment(db *sql.DB, id int) (comment Comment, err error) {
	err = db.QueryRow("SELECT content, author, post_id FROM comments WHERE id = $1",
		id).Scan(&comment.Content, &comment.Author, &comment.Post_Id)
	return
}

func (comment *Comment) Update(db *sql.DB, newComment string) (err error) {
	_, err = db.Exec("UPDATE comments SET content = $1 WHERE id = $2",
		newComment, comment.Id)
	comment.Content = newComment
	return
}

// Returns an array of "comment" structures from the database by the post ID.
func GetAllCommentsByPost(db *sql.DB, post_id int) ([]Comment, error) {
	var result []Comment

	statement := "SELECT id, content, author FROM comments WHERE post_id = $1"

	rows, err := db.Query(statement, post_id)
	if err != nil {
		return nil, fmt.Errorf("get comments by post: %w", err)
	}

	for rows.Next() {
		var comment Comment

		err = rows.Scan(&comment.Id, &comment.Content, &comment.Author)
		if err != nil {
			return nil, fmt.Errorf("get comments by post: %w", err)
		}

		comment.Post_Id = post_id
		result = append(result, comment)
	}

	return result, nil
}
