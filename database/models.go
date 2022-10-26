package database

type Post struct {
	Id         int
	Title      string
	Message    string
	Author     string
	User_Id    int
	Like       int
	Dislike    int
	Categories []string
}

type User struct {
	Id        int
	Firstname string
	Lastname  string
	Nickname  string
	Email     string
	Posts     []Post
}

type Comment struct {
	Id      int
	Content string
	Author  string
	Like    int
	Dislike int
	Post    *Post
}

type Category struct {
	Id   int
	Name string
	Post []*Post
}
