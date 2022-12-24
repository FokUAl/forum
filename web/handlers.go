package web

import (
	"fmt"
	"forumAA/database"
	"html/template"
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
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	t, err := template.ParseFiles("./ui/template/home.html")
	if err != nil {
		app.errorLog.Printf("home: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	user := app.checkUser(w, r)
	posts, err := database.GetAllPost(app.database)
	if err != nil {
		app.errorLog.Printf("home: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	catFilter := r.FormValue("catFilter")
	if catFilter != "" {
		if r.Method != http.MethodPost {
			app.errorLog.Printf("likeComment: method not allowed\n")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		posts, err = database.GetPostByCategory(app.database, catFilter)
		if err != nil {
			app.errorLog.Printf("home: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
	}

	new_info := info{
		User:  user,
		Posts: posts,
	}

	err = t.Execute(w, new_info)
	if err != nil {
		app.errorLog.Printf("home: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}

func (app *application) profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorLog.Printf("profile: Method not allowed\n")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	user := app.checkUser(w, r)
	if user.Id == 0 {
		app.errorLog.Printf("profile: unauthorized user\n")
		http.Error(w, http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
		return
	}

	t, err := template.ParseFiles("./ui/template/profile.html")
	if err != nil {
		app.errorLog.Printf("profile: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	createdPosts, err := database.GetPostsByUser(app.database, user.Id)
	if err != nil {
		app.errorLog.Printf("profile: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	likedPost, err := database.GetPostsLikedByUser(app.database, user.Nickname)
	if err != nil {
		app.errorLog.Printf("profile: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	new_info := info{
		User:       user,
		Posts:      createdPosts,
		LikedPosts: likedPost,
	}
	err = t.Execute(w, new_info)
	if err != nil {
		app.errorLog.Printf("profile: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) likeComment(w http.ResponseWriter, r *http.Request) {
	user := app.checkUser(w, r)
	if user.Id == 0 {
		app.errorLog.Printf("likeComment: unauthorized user\n")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	comment_id_str := strings.TrimPrefix(r.URL.Path, "/comment/like/")
	comment_id, err := strconv.Atoi(comment_id_str)
	if err != nil {
		app.errorLog.Printf("likeComment: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		app.errorLog.Printf("likeComment: method not allowed\n")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	like, err := database.GetCommentLikeByUser(app.database, user.Nickname, comment_id)
	if err != nil {
		app.infoLog.Printf("likeComment: %s\n", err.Error())
		// http.Error(w, http.StatusText(http.StatusInternalServerError),
		// 	http.StatusInternalServerError)
		// return
	}

	comment, err := database.GetComment(app.database, comment_id)
	if err != nil {
		app.errorLog.Printf("likeComment: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.errorLog.Printf("likeComment: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	likeValue, err := strconv.Atoi(r.FormValue("commentLikeBtn"))
	if err != nil {
		app.errorLog.Printf("likeComment: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
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
		app.errorLog.Printf("likeComment: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment.Post_Id), http.StatusSeeOther)
}
