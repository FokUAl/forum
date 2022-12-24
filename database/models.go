package database

type Post struct {
	Id         int
	Title      string
	Message    string
	Author     string
	User_Id    int
	Like       int
	Categories []string
}

type User struct {
	Id        int
	Firstname string
	Lastname  string
	Nickname  string
	Password  string
	Email     string
	Posts     []Post
}

type Comment struct {
	Id      int
	Content string
	Author  string
	Like    int
	Post_Id int
}

type Category struct {
	Id   int
	Name string
	Post []*Post
}

type Like struct {
	Id       int
	Nickname string
	Value    int
	Elem_Id  int
}
