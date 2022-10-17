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

func DeleteCategory(category_id int) (err error) {
	_, err = database.Exec("DELETE FROM categories WHERE id = $1", category_id)
	return
}
