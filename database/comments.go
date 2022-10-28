package database

import (
	"database/sql"
	"errors"
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

func GetComment(db *sql.DB, id int) (err error) {
	comment := Comment{}
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
