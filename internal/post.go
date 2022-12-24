package internal

import (
	"database/sql"
	"fmt"
	"forumAA/database"
	"net/http"
	"strings"
)

func CreatePost(db *sql.DB, user database.User, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	title := r.FormValue("postTitle")
	message := r.FormValue("postMessage")
	categories := r.Form["postCat"]

	title = strings.Trim(title, " ")
	message = strings.Trim(message, " ")

	if title == "" || message == "" {
		return fmt.Errorf("Title or message is empty")
	}

	for i := 0; i < len(categories); i++ {
		categories[i] = strings.Trim(categories[i], " ")
		if categories[i] == "" {
			return fmt.Errorf("Category is empty")
		}
	}

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

func CountPostLikes(db *sql.DB, post_id int) (int, error) {
	likes, err := database.GetLikeByPost(db, post_id)
	if err != nil {
		return 0, fmt.Errorf("CountPostLikes: %w", err)
	}

	result := 0
	for _, like := range likes {
		result += like.Value
	}

	return result, nil
}
