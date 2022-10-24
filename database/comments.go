package database

import "errors"

type Comment struct {
	Id      int
	Content string
	Author  string
	Post    *Post
}

func (comment *Comment) Create() (err error) {
	if comment.Post == nil {
		err = errors.New("Post not found")
		return
	}
	statement := "INSERTS INTO comments (content, author, post_id)" +
		"values ($1, $2, $3)"
	err = database.QueryRow(statement, comment.Content, comment.Author, comment.Post.Id).Scan(&comment.Id)
	return
}

func (comment *Comment) Delete() (err error) {
	_, err = database.Exec("DELETE FROM comments where id = $1", comment.Id)
	return
}

func (comment *Comment) Update() (err error) {
	_, err = database.Exec("UPDATE comments SET content = $2, author = $3 WHERE id = $1",
		comment.Id, comment.Content, comment.Author)
	return
}
