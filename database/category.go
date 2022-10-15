package database

type Category struct {
	Id   int
	Name string
	Post []string
}

func CreateCategory(category_name string, post Post) (err error) {
	_, err = database.Exec("INSERT INTO categories (name, post) values ($1, $2)", category_name, post)
	return
}
