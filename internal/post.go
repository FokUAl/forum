package internal

import (
	"database/sql"
	"fmt"
	"forumAA/database"
	"net/http"
)

func CreatePost(db *sql.DB, user database.User, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	title := r.FormValue("postTitle")
	message := r.FormValue("postMessage")
	categories := r.Form["postCat"]

	post := database.Post{
		Title:      title,
		Message:    message,
		Author:     user.Nickname,
		User_Id:    user.Id,
		Categories: categories,
	}

	err = post.Create(db)
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
