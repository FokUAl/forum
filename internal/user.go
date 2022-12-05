package internal

import (
	"database/sql"
	"errors"
	"forumAA/database"

	"golang.org/x/crypto/bcrypt"
)

func Registration(db *sql.DB, new_user database.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(new_user.Password), 8)
	new_user.Password = string(hashedPassword)

	err = new_user.Create(db)

	return err
}

func Login(db *sql.DB, user_nick string, user_passw string) error {
	userInfo, err := database.GetUser(db, user_nick)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(user_passw))
	if err != nil {
		err = errors.New("Login: password or login invalid")
	}

	return err
}
