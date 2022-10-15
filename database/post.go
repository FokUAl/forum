package database

type Post struct {
	Id         int
	Message    string
	Author     string
	rating     int
	Categories []string
}

func (post *Post) Create() (err error) {
	statement := "INSERT INTO posts (message, author) values ($1, $2) returning id"
	stmt, err := database.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(post.Message, post.Author).Scan(&post.Id)
	return
}

func (post *Post) Delete() (err error) {
	_, err = database.Exec("DELETE FROM posts where id = $1", post.Id)
	return
}

func GetPost(id int) (post Post, err error) {
	post = Post{}
	err = database.QueryRow("select id, content, author from posts where id = $1",
		id).Scan(&post.Id, &post.Message, &post.Author)
	return
}

func (post *Post) Update() (err error) {
	_, err = database.Exec("update posts set content = $2, author = $3 where id = $1",
		post.Id, post.Message, post.Author)
	return
}
