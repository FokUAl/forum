package web

import (
	"fmt"
	"forumAA/database"
	"forumAA/internal"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type info struct {
	Posts    []database.Post
	User     database.User
	Post     database.Post
	Comments []database.Comment
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

func (app *application) post(w http.ResponseWriter, r *http.Request) {
	post_id_str := strings.TrimPrefix(r.URL.Path, "/post/")
	post_id, err := strconv.ParseInt(post_id_str, 10, 32)
	if err != nil {
		http.Error(w, "post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user := app.checkUser(w, r)
	post, err := database.GetPost(app.database, int(post_id))
	if err != nil {
		http.Error(w, "post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	post_likes, err := internal.CountLikes(app.database, post.Id)
	if err != nil {
		http.Error(w, "post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	post.Like = post_likes

	switch r.Method {
	case http.MethodGet:
		comments, err := database.GetAllCommentsByPost(app.database, post.Id)
		if err != nil {
			http.Error(w, "post: "+err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.ParseFiles("./ui/template/post.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "File not found: post.html", 500)
			return
		}

		new_info := info{
			User:     user,
			Post:     post,
			Comments: comments,
		}

		err = t.Execute(w, new_info)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

		comment_content := r.FormValue("comment")
		comment := database.Comment{
			Content: comment_content,
			Author:  user.Nickname,
			Post:    &post,
		}

		comment.Create(app.database)
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	}
}

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		http.NotFound(w, r)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("./ui/template/signUp.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "File not found: signUp.html", 500)
			return
		}

		err = t.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

		user := database.User{
			Firstname: r.FormValue("firstName"),
			Lastname:  r.FormValue("lastName"),
			Email:     r.FormValue("email"),
			Nickname:  r.FormValue("nickname"),
			Password:  r.FormValue("password"),
		}

		err = internal.Registration(app.database, user)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) signIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin" {
		http.NotFound(w, r)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("./ui/template/signIn.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "File not found: signin.html", 500)
			return
		}

		err = t.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

		nick := r.FormValue("nickname")
		password := r.FormValue("password")

		err = internal.Login(app.database, nick, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		sessionToken, err := uuid.NewV4()
		if err != nil {
			log.Fatalf("failed to generate UUID: %v", err)
		}

		expiresAt := time.Now().Add(12 * time.Hour)

		database.CreateSession(app.database, nick, sessionToken.String(), expiresAt)

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken.String(),
			Expires: expiresAt,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
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

	err = t.Execute(w, user)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) logOut(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return

	}
	sessionToken := c.Value
	database.DeleteSession(app.database, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) likeComment(w http.ResponseWriter, r *http.Request) {
	user := app.checkUser(w, r)
	if user.Id == 0 {
		http.Error(w, "likeComment: unauthorized user", http.StatusUnauthorized)
		return
	}

	comment_id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/like/"))
	if err != nil {
		http.Error(w, "likeComment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "likeComment: method not allowed", http.StatusMethodNotAllowed)
		return
	}

	likes, err := database.GetLikeByComment(app.database, comment_id)
	if err != nil {
		http.Error(w, "likeComment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	likeValue, err := strconv.Atoi(r.FormValue("likeBtn"))
	if err != nil {
		http.Error(w, "likeComment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, like := range likes {
		if like.Nickname == user.Nickname && like.Value != likeValue {
			database.UpdateCommentLike(app.database, likeValue, user.Nickname, comment_id)
			http.Redirect(w, r, fmt.Sprintf("/post/%d", comment_id), http.StatusSeeOther)
		} else if like.Nickname == user.Nickname {
			database.UpdateCommentLike(app.database, 0, user.Nickname, comment_id)
			http.Redirect(w, r, fmt.Sprintf("/post/%d", comment_id), http.StatusSeeOther)
		}
	}

	if likeValue > 0 {
		err = database.CreateCommentLike(app.database, user.Nickname, 1, comment_id)
	} else {
		err = database.CreateCommentLike(app.database, user.Nickname, 1, comment_id)
	}

	if err != nil {
		http.Error(w, "likeComment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment_id), http.StatusSeeOther)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	user := app.checkUser(w, r)
	if user.Id == 0 {
		http.Error(w, "likeComment: unauthorized user", http.StatusUnauthorized)
		return
	}

	if r.URL.Path != "/create-post" {
		http.Error(w, "createPost: Not Found", http.StatusNotFound)
	}

	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("./ui/template/createPost.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "File not found: signin.html", 500)
			return
		}

		err = t.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	case http.MethodPost:
		err := internal.CreatePost(app.database, user, r)
		if err != nil {
			http.Error(w, "createPost: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
