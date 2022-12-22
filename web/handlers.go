package web

import (
	"fmt"
	"forumAA/database"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type info struct {
	Posts      []database.Post
	User       database.User
	Post       database.Post
	Comments   []database.Comment
	LikedPosts []database.Post
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	t, err := template.ParseFiles("./ui/template/home.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "File not found: home.html", 500)
		return
	}

	user := app.checkUser(w, r)
	posts, err := database.GetAllPost(app.database)
	if err != nil {
		http.Error(w, "home: get all post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	new_info := info{
		User:  user,
		Posts: posts,
	}
	err = t.Execute(w, new_info)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	user := app.checkUser(w, r)
	if user.Id == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	t, err := template.ParseFiles("./ui/template/profile.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "File not found: profile.html", 500)
		return
	}

	createdPosts, err := database.GetPostsByUser(app.database, user.Id)
	if err != nil {
		http.Error(w, "profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	likedPost, err := database.GetPostsLikedByUser(app.database, user.Nickname)
	if err != nil {
		http.Error(w, "profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	new_info := info{
		User:       user,
		Posts:      createdPosts,
		LikedPosts: likedPost,
	}
	err = t.Execute(w, new_info)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) likeComment(w http.ResponseWriter, r *http.Request) {
	user := app.checkUser(w, r)
	if user.Id == 0 {
		http.Error(w, "likeComment: unauthorized user", http.StatusUnauthorized)
		return
	}

	comment_id_str := strings.TrimPrefix(r.URL.Path, "/comment/like/")
	comment_id, err := strconv.Atoi(comment_id_str)
	if err != nil {
		http.Error(w, "likeComment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "likeComment: method not allowed", http.StatusMethodNotAllowed)
		return
	}

	like, err := database.GetCommentLikeByUser(app.database, user.Nickname, comment_id)

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	likeValue, err := strconv.Atoi(r.FormValue("commentLikeBtn"))
	if err != nil {
		http.Error(w, "likeComment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if like.Id == 0 && likeValue > 0 {
		err = database.CreateCommentLike(app.database, user.Nickname, 1, comment_id)
	} else if like.Id == 0 && likeValue < 0 {
		err = database.CreateCommentLike(app.database, user.Nickname, -1, comment_id)
	} else if likeValue == like.Value {
		err = database.DeleteCommentLike(app.database, like.Id)
	} else if likeValue != like.Value {
		err = database.UpdateCommentLike(app.database, likeValue, user.Nickname, comment_id)
	}

	if err != nil {
		http.Error(w, "likeComment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment_id), http.StatusSeeOther)
}
