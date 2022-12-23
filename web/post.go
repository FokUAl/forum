package web

import (
	"fmt"
	"forumAA/database"
	"forumAA/internal"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) post(w http.ResponseWriter, r *http.Request) {
	post_id_str := strings.TrimPrefix(r.URL.Path, "/post/")
	post_id, err := strconv.ParseInt(post_id_str, 10, 32)
	if err != nil {
		app.errorLog.Printf("post: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	user := app.checkUser(w, r)
	post, err := database.GetPost(app.database, int(post_id))
	if err != nil {
		app.errorLog.Printf("post: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	post_likes, err := internal.CountPostLikes(app.database, post.Id)
	if err != nil {
		app.errorLog.Printf("post: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	post.Like = post_likes

	switch r.Method {
	case http.MethodGet:
		comments, err := database.GetAllCommentsByPost(app.database, post.Id)
		if err != nil {
			app.errorLog.Printf("post: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		for i := 0; i < len(comments); i++ {
			comment_likes, err := internal.CountCommentLikes(app.database, comments[i].Id)
			if err != nil {
				app.errorLog.Printf("post: %s\n", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			comments[i].Like = comment_likes
		}

		t, err := template.ParseFiles("./ui/template/post.html")
		if err != nil {
			app.errorLog.Printf("post: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		new_info := info{
			User:     user,
			Post:     post,
			Comments: comments,
		}

		err = t.Execute(w, new_info)
		if err != nil {
			app.errorLog.Printf("post: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			app.errorLog.Printf("post: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}

		comment_content := r.FormValue("comment")
		comment := database.Comment{
			Content: comment_content,
			Author:  user.Nickname,
			Post:    &post,
		}

		err = comment.Create(app.database)
		if err != nil {
			app.errorLog.Printf("post: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	}
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	user := app.checkUser(w, r)
	if user.Id == 0 {
		app.errorLog.Printf("createPost: unauthorized user\n")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if r.URL.Path != "/create-post" {
		app.errorLog.Printf("createPost: invalid path %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("./ui/template/createPost.html")
		if err != nil {
			app.errorLog.Printf("createPost: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, user)
		if err != nil {
			app.errorLog.Printf("createPost: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	case http.MethodPost:
		err := internal.CreatePost(app.database, user, r)
		if err != nil {
			app.errorLog.Printf("createPost: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *application) likePost(w http.ResponseWriter, r *http.Request) {
	user := app.checkUser(w, r)
	if user.Id == 0 {
		app.errorLog.Printf("likePost: unauthorized user\n")
		http.Error(w, http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
		return
	}

	post_id_str := strings.TrimPrefix(r.URL.Path, "/post/like/")
	post_id, err := strconv.Atoi(post_id_str)
	if err != nil {
		app.errorLog.Printf("likePost: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		app.errorLog.Printf("likePost: Method not allowed\n")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
		return
	}

	like, err := database.GetPostLikeByUser(app.database, user.Nickname, post_id)
	if err != nil {
		app.infoLog.Printf("likePost: %s\n", err.Error())
		// http.Error(w, http.StatusText(http.StatusInternalServerError),
		// 	http.StatusInternalServerError)
		// return
	}

	err = r.ParseForm()
	if err != nil {
		app.errorLog.Printf("likePost: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	likeValue, err := strconv.Atoi(r.FormValue("postLikeBtn"))
	if err != nil {
		app.errorLog.Printf("likePost: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if like.Id == 0 && likeValue > 0 {
		err = database.CreatePostLike(app.database, user.Nickname, 1, post_id)
	} else if like.Id == 0 && likeValue < 0 {
		err = database.CreatePostLike(app.database, user.Nickname, -1, post_id)
	} else if likeValue == like.Value {
		err = database.DeletePostLikes(app.database, like.Id)
	} else if likeValue != like.Value {
		err = database.UpdatePostLike(app.database, likeValue, user.Nickname, post_id)
	}

	if err != nil {
		app.errorLog.Printf("likePost: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", post_id), http.StatusSeeOther)
}
