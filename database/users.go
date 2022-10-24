package database

type User struct {
	Id        int
	Firstname string
	Lastname  string
	Nickname  string
	Email     string
	Posts     []Post
}

func (user *User) Create() (err error) {
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

func (user *User) Delete() (err error) {
	_, err = database.Exec("DELETE FROM posts where id = $1", user.Id)
	return
}

func GetUser(id int) (user User, err error) {
	user = User{}
	err = database.QueryRow("SELECT id, firstname, lastname, nickname, email FROM users WHERE id = $1",
		id).Scan(&user.Firstname, &user.Lastname, &user.Nickname, &user.Email)
	return
}

func (user *User) Update() (err error) {
	_, err = database.Exec("UPDATE users SET firstname = $2, lastname = $3 WHERE id = $1",
		user.Id, user.Firstname, user.Lastname)
	return
}
