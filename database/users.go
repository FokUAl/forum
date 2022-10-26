package database

import (
	"database/sql"
)

func (user *User) Create(db *sql.DB) (err error) {
	statement := "INSERT INTO users (firstname, lastname, nickname, email) " +
		"VALUES ($1, $2, $3, $4) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(user.Firstname, user.Lastname, user.Nickname, user.Email).Scan(&user.Id)
	return
}

func (user *User) Delete(db *sql.DB) (err error) {
	_, err = db.Exec("DELETE FROM users where id = $1", user.Id)
	return
}

func GetUser(db *sql.DB, id int) (user User, err error) {
	user = User{}
	err = db.QueryRow("SELECT id, firstname, lastname, nickname, email FROM users WHERE id = $1",
		id).Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Nickname, &user.Email)
	return
}

func (user *User) Update(db *sql.DB, newFirstname, newLastname string) (err error) {
	_, err = db.Exec("UPDATE users SET firstname = $1, lastname = $2 WHERE id = $3",
		newFirstname, newLastname, user.Id)

	user.Firstname = newFirstname
	user.Lastname = newLastname

	return
}
