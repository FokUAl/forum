package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (comment *Comment) Create(db *sql.DB) (err error) {
	if comment.Post == nil {
		err = errors.New("Post not found")
		return
	}
	statement := "INSERT INTO comments (content, author, post_id)" +
		"values ($1, $2, $3) returning id"
	err = db.QueryRow(statement, comment.Content, comment.Author, comment.Post.Id).Scan(&comment.Id)
	return
}

func (comment *Comment) Delete(db *sql.DB) (err error) {
	_, err = db.Exec("DELETE FROM comments where id = $1", comment.Id)
	return
}

func GetComment(db *sql.DB, id int) (comment Comment, err error) {
	err = db.QueryRow("SELECT content, author, like, dislike FROM comments WHERE id = $1",
		id).Scan(&comment.Content, comment.Author, comment.Like, comment.Dislike)
	return
}

func (comment *Comment) Update(db *sql.DB, newComment string) (err error) {
	_, err = db.Exec("UPDATE comments SET content = $1 WHERE id = $2",
		newComment, comment.Id)
	comment.Content = newComment
	return
}

func GetAllCommentsByPost(db *sql.DB, post_id int) ([]Comment, error) {
	var result []Comment

	statement := "SELECT * FROM comments WHERE post_id = $1"

	rows, err := db.Query(statement, post_id)
	if err != nil {
		return nil, fmt.Errorf("get comments by post: %w", err)
	}

	for rows.Next() {
		var comment Comment

		post, err := GetPost(db, post_id)
		if err != nil {
			return nil, fmt.Errorf("get comments by post: %w", err)
		}

		var _id int
		err = rows.Scan(&comment.Id, &comment.Content, &comment.Author, &comment.Like, &comment.Dislike, &_id)
		if err != nil {
			return nil, fmt.Errorf("get comments by post: %w", err)
		}

		comment.Post = &post
		result = append(result, comment)
	}

	return result, nil
}
