package database

import (
	"database/sql"
	"time"
)

func GetUserByToken(db *sql.DB, token string) (User, error) {
	nick := ""
	err := db.QueryRow("SELECT nickname FROM sessions WHERE session_token = $1",
		token).Scan(&nick)
	if err != nil {
		return User{}, err
	}

	user, err := GetUser(db, nick)

	return user, err
}

func GetExpiryByToken(db *sql.DB, token string) (time.Time, error) {
	result := time.Time{}

	err := db.QueryRow("SELECT expiry FROM sessions WHERE session_token = $1",
		token).Scan(&result)

	return result, err
}

func DeleteSession(db *sql.DB, token string) error {
	_, err := db.Exec("DELETE FROM sessions where session_token = $1", token)

	return err
}

func CreateSession(db *sql.DB, nick string, token string, expiry time.Time) error {
	statement := "INSERT INTO sessions (nickname, session_token, expiry) VALUES ($1, $2, $3) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()
	id := 0
	err = stmt.QueryRow(nick, token, expiry).Scan(&id)
	return err
}

func UpdateSession(db *sql.DB, nick string, new_token string, expiry time.Time) error {
	_, err := db.Exec("UPDATE sessions SET session_token = $1 WHERE nickname = $2", new_token, nick)

	return err
}

func IsSessionExist(db *sql.DB, nick string) (bool, error) {
	id := 0
	err := db.QueryRow("SELECT id FROM sessions WHERE nickname = $1", nick).Scan(&id)
	if err != nil {
		return false, err
	}

	return id != 0, nil
}
