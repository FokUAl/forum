package database

import (
	"database/sql"
)

func (post *Post) Create(database *sql.DB) (err error) {
	statement := "INSERT INTO posts (title, message, author, user_id) VALUES ($1, $2, $3, $4) returning id"
	stmt, err := database.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(post.Title, post.Message, post.Author, post.User_Id).Scan(&post.Id)
	return
}

func (post *Post) Delete(database *sql.DB) (err error) {
	_, err = database.Exec("DELETE FROM posts where id = $1", post.Id)
	return
}

func GetPost(database *sql.DB, id int) (post Post, err error) {
	post = Post{}
	err = database.QueryRow("SELECT id, title, message, author FROM posts WHERE id = $1",
		id).Scan(&post.Id, &post.Title, &post.Message, &post.Author)
	return
}

func (post *Post) Update(database *sql.DB, newTitle, newMessage string) (err error) {
	_, err = database.Exec("UPDATE posts SET message = $1, title = $2 WHERE id = $3",
		newMessage, newTitle, post.Id)

	post.Title = newTitle
	post.Message = newMessage

	return
}
