package database

import (
	"database/sql"
)

func (user *User) Create(db *sql.DB) (err error) {
	statement := "INSERT INTO users (firstname, lastname, nickname, password, email) " +
		"VALUES ($1, $2, $3, $4, $5) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(user.Firstname, user.Lastname, user.Nickname, user.Password, user.Email).Scan(&user.Id)
	return
}

func (user *User) Delete(db *sql.DB) (err error) {
	_, err = db.Exec("DELETE FROM users where id = $1", user.Id)
	return
}

func GetUser(db *sql.DB, nickname string) (user User, err error) {
	user = User{}
	err = db.QueryRow("SELECT id, firstname, lastname, nickname, password, email FROM users WHERE nickname = $1",
		nickname).Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Nickname, &user.Password, &user.Email)
	return
}

func (user *User) Update(db *sql.DB, newFirstname, newLastname string) (err error) {
	_, err = db.Exec("UPDATE users SET firstname = $1, lastname = $2 WHERE id = $3",
		newFirstname, newLastname, user.Id)

	user.Firstname = newFirstname
	user.Lastname = newLastname

	return
}

func (user *User) ChangePassword(db *sql.DB, new_password string) error {
	_, err := db.Exec("UPDATE users SET password = $1 WHERE id = $2",
		new_password, user.Id)

	user.Password = new_password

	return err
}

func IsNicknameExist(db *sql.DB, nick string) (bool, error) {
	var existing_id int
	err := db.QueryRow("SELECT id FROM users WHERE nickname = $1",
		nick).Scan(&existing_id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
