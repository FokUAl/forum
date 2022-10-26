package database

import (
	"database/sql"
)

func (user *User) Create(database *sql.DB) (err error) {
	statement := "INSERT INTO users (firstname, lastname, nickname, email) " +
		"VALUES ($1, $2, $3, $4) returning id"
	stmt, err := database.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(user.Firstname, user.Lastname, user.Nickname, user.Email).Scan(&user.Id)
	return
}

func (user *User) Delete(database *sql.DB) (err error) {
	_, err = database.Exec("DELETE FROM users where id = $1", user.Id)
	return
}

func GetUser(database *sql.DB, id int) (user User, err error) {
	user = User{}
	err = database.QueryRow("SELECT id, firstname, lastname, nickname, email FROM users WHERE id = $1",
		id).Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Nickname, &user.Email)
	return
}

func (user *User) Update(database *sql.DB, newFirstname, newLastname string) (err error) {
	_, err = database.Exec("UPDATE users SET firstname = $1, lastname = $2 WHERE id = $3",
		newFirstname, newLastname, user.Id)

	user.Firstname = newFirstname
	user.Lastname = newLastname

	return
}
