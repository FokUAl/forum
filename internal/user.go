package internal

import (
	"database/sql"
	"errors"
	"forumAA/database"
)

func Registration(db *sql.DB, new_user database.User) error {
	// TO DO
	// Reading data from front

	user := new_user

	err := user.Create(db)

	return err
}

func Login(db *sql.DB, user_nick string, user_passw string) error {
	// TO DO
	// Reading nickname and password from front

	userInfo, err := database.GetUser(db, user_nick)
	if err != nil {
		return err
	}

	if user_passw != userInfo.Password {
		err = errors.New("Login: password or login invalid")
	}

	return err
}
