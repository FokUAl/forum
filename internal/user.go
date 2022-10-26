package internal

import (
	"database/sql"
	"errors"
	"forumAA/database"
)

func Registration(db *sql.DB) error {
	// TO DO
	// Reading data from front

	user := database.User{}

	err := user.Create(db)

	return err
}

func Login(db *sql.DB) error {
	// TO DO
	// Reading nickname and password from front
	nick := ""
	pass := ""

	userInfo, err := database.GetUser(db, nick)
	if err != nil {
		return err
	}

	if pass != userInfo.Password {
		err = errors.New("Login: password or login invalid")
	}

	return err
}
