package database

import (
	"database/sql"
	"fmt"
)

func GetLikeByPost(db *sql.DB, post_id int) ([]Like, error) {
	result := []Like{}
	statement := "SELECT id, nickname, like FROM post_likes WHERE post_id = $1"
	rows, err := db.Query(statement, post_id)
	if err != nil {
		return nil, fmt.Errorf("GetPostLike: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var like Like

		err := rows.Scan(&like.Id, &like.Nickname, &like.Value)
		if err != nil {
			return nil, fmt.Errorf("GetPostLike: %w", err)
		}

		like.Elem_Id = post_id

		result = append(result, like)
	}

	return result, nil
}

func GetPostLikeByUser(db *sql.DB, nickname string, post_id int) (Like, error) {
	result := Like{}

	err := db.QueryRow("SELECT id, like FROM post_likes "+
		"WHERE nickname = $1 AND post_id = $2", nickname,
		post_id).Scan(&result.Id, &result.Value)
	if err != nil {
		return result, fmt.Errorf("GetPostLikeByUser: %w", err)
	}

	result.Nickname = nickname
	result.Elem_Id = post_id

	return result, nil
}

func UpdatePostLike(db *sql.DB, new_value int, nickname string, post_id int) error {
	_, err := db.Exec("UPDATE post_likes SET like = $1 "+
		"WHERE nickname = $2 AND post_id = $3",
		new_value, nickname, post_id)

	return err
}

func CreatePostLike(db *sql.DB, nickname string, like int, post_id int) error {
	statement := "INSERT INTO post_likes (nickname, like, post_id) VALUES ($1, $2, $3)"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()
	err = stmt.QueryRow(nickname, like, post_id).Scan()
	return err
}

func GetCommentLikeByUser(db *sql.DB, nickname string, comment_id int) (Like, error) {
	result := Like{}

	err := db.QueryRow("SELECT id, like FROM comment_likes "+
		"WHERE nickname = $1 AND comment_id = $2", nickname,
		comment_id).Scan(&result.Id, &result.Value)
	if err != nil {
		return result, fmt.Errorf("GetPostLikeByUser: %w", err)
	}

	result.Nickname = nickname
	result.Elem_Id = comment_id

	return result, nil
}

func GetLikeByComment(db *sql.DB, comment_id int) ([]Like, error) {
	result := []Like{}

	statement := "SELECT id, nickname, like FROM comment_likes WHERE comment_id = $1"
	rows, err := db.Query(statement, comment_id)
	if err != nil {
		return nil, fmt.Errorf("GetLikeByComment: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var like Like

		err := rows.Scan(&like.Id, &like.Nickname, &like.Value)
		if err != nil {
			return nil, fmt.Errorf("GetPostLike: %w", err)
		}

		like.Elem_Id = comment_id

		result = append(result, like)
	}

	return result, nil
}

func UpdateCommentLike(db *sql.DB, new_value int, nickname string, comment_id int) error {
	_, err := db.Exec("UPDATE comment_likes SET like = $1 "+
		"WHERE nickname = $2 AND comment_id = $3",
		new_value, nickname, comment_id)

	return err
}

func CreateCommentLike(db *sql.DB, nickname string, like int, comment_id int) error {
	statement := "INSERT INTO comment_likes (nickname, like, comment_id) VALUES ($1, $2, $3) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()
	id := 0
	err = stmt.QueryRow(nickname, like, comment_id).Scan(&id)
	return err
}
