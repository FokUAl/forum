package database

import (
	"database/sql"
	"fmt"
)

func GetPostLike(db *sql.DB, post_id int) ([]Like, error) {
	result := []Like{}

	rows, err := db.Query("SELECT id, nickname, like FROM post_likes WHERE post_id = $1")
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

// func UpdatePostLike(db *sql.DB, new_value int, post_id int) error {
// 	_, err := db.Exec("UPDATE post_likes SET like = $1 WHERE post_id = $2",
// 		new_value, post_id)

// 	return err
// }

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
