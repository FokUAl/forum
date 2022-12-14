package web

import (
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
		w.WriteHeader(http.StatusNotFound)
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
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	post_id_str := strings.TrimPrefix(r.URL.Path, "/post/")
	post_id, err := strconv.ParseInt(post_id_str, 10, 32)
	if err != nil {
		http.Error(w, "post: Internal Error", http.StatusInternalServerError)
		return
	}

	post, err := database.GetPost(app.database, int(post_id))
	if err != nil {
		http.Error(w, "post: Internal Error", http.StatusInternalServerError)
		return
	}

	comments, err := database.GetAllCommentsByPost(app.database, post.Id)

	t, err := template.ParseFiles("./ui/template/post.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "File not found: post.html", 500)
		return
	}

	user := app.checkUser(w, r)

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
