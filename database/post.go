package database

import (
	"database/sql"
	"fmt"
)

// Creates row in posts table by post structure.
func (post *Post) Create(db *sql.DB) (err error) {
	statement := "INSERT INTO posts (title, message, author, user_id) VALUES ($1, $2, $3, $4) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}

	err = stmt.QueryRow(post.Title, post.Message, post.Author, post.User_Id).Scan(&post.Id)
	if err != nil {
		return err
	}

	stmt.Close()

	err = CreateCategories(db, post.Categories, post.Id)

	return
}

// Deletes row from posts table by post structure.
func (post *Post) Delete(db *sql.DB) (err error) {
	_, err = db.Exec("DELETE FROM posts where id = $1", post.Id)
	return
}

// Gets posts structure by post id.
func GetPost(db *sql.DB, id int) (post Post, err error) {
	post = Post{}
	err = db.QueryRow("SELECT id, title, message, author FROM posts WHERE id = $1",
		id).Scan(&post.Id, &post.Title, &post.Message, &post.Author)

	if err != nil {
		return post, err
	}
	categories, err := GetAllCategoryByPost(db, id)
	if err != nil {
		return Post{}, fmt.Errorf("get all posts: %w", err)
	}

	post.Categories = categories
	return
}

func (post *Post) Update(db *sql.DB, newTitle, newMessage string) (err error) {
	_, err = db.Exec("UPDATE posts SET message = $1, title = $2 WHERE id = $3",
		newMessage, newTitle, post.Id)

	post.Title = newTitle
	post.Message = newMessage

	return
}

func GetAllPost(db *sql.DB) ([]Post, error) {
	statement := "SELECT * FROM posts"

	rows, err := db.Query(statement)
	if err != nil {
		return nil, fmt.Errorf("get all posts: %w", err)
	}

	defer rows.Close()

	var result []Post

	for rows.Next() {
		var post Post
		var user_id int
		err := rows.Scan(&post.Id, &post.Title, &post.Message,
			&post.Author, &user_id)
		if err != nil {
			return nil, fmt.Errorf("get all posts: %w", err)
		}

		categories, err := GetAllCategoryByPost(db, post.Id)
		if err != nil {
			return nil, fmt.Errorf("get all posts: %w", err)
		}

		post.Categories = categories

		result = append(result, post)
	}

	return result, nil
}

func GetPostByCategory(db *sql.DB, category string) ([]Post, error) {
	var result []Post

	statement := "SELECT id, title, message, author FROM posts WHERE id IN (SELECT post_id FROM categories WHERE name = $1)"
	rows, err := db.Query(statement, category)
	if err != nil {
		return nil, fmt.Errorf("get post by category: %w", err)
	}

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Title, &post.Message,
			&post.Author)
		if err != nil {
			return nil, fmt.Errorf("get posts by category: %w", err)
		}

		categories, err := GetAllCategoryByPost(db, post.Id)
		if err != nil {
			return nil, fmt.Errorf("get posts by category: %w", err)
		}

		post.Categories = categories

		result = append(result, post)
	}

	return result, nil
}

func GetPostsByUser(db *sql.DB, user_id int) ([]Post, error) {
	var result []Post

	statement := "SELECT id, title, message, author FROM posts WHERE user_id = $1"

	rows, err := db.Query(statement, user_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Title, &post.Message,
			&post.Author)
		if err != nil {
			return nil, fmt.Errorf("get posts by user: %w", err)
		}

		categories, err := GetAllCategoryByPost(db, post.Id)
		if err != nil {
			return nil, fmt.Errorf("get posts by user: %w", err)
		}

		post.Categories = categories

		result = append(result, post)
	}

	return result, nil
}

func GetPostsLikedByUser(db *sql.DB, nickname string) ([]Post, error) {
	var result []Post

	likes, err := GetPostLikesByUser(db, nickname)
	if err != nil {
		return nil, fmt.Errorf("get post likes by user: %w", err)
	}

	for _, like := range likes {
		if like.Nickname == nickname {
			post, err := GetPost(db, like.Elem_Id)
			if err != nil {
				return nil, fmt.Errorf("get post likes by user: %w", err)
			}

			result = append(result, post)
		}
	}

	return result, nil
}
