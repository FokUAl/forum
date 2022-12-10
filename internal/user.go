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

func NewPassword(db *sql.DB, user_nick string, old_password string, new_password string) error {
	userInfo, err := database.GetUser(db, user_nick)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(old_password))
	if err != nil {
		err = errors.New("NewPassword: password is invalid")
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(new_password), 8)
	if err != nil {
		return err
	}

	userInfo.ChangePassword(db, string(hashedPassword))

	return err
}
