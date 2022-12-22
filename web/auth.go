package web

import (
	"forumAA/database"
	"forumAA/internal"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

type notification struct {
	Content string
	Exist   bool
}

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sign-up" {
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

		err = t.Execute(w, app.notice)
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

		content, ok := internal.CheckInput(app.database, user)

		if !ok {
			app.notice = notification{
				Content: content,
				Exist:   true,
			}
			http.Redirect(w, r, "/sign-up", http.StatusSeeOther)
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
	if r.URL.Path != "/sign-in" {
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
			return
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
