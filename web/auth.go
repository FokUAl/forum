package web

import (
	"forumAA/database"
	"forumAA/internal"
	"html/template"
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
		app.errorLog.Printf("signUp: invalid path: %s\n", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("./ui/template/signUp.html")
		if err != nil {
			app.errorLog.Printf("signUp: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, app.notice)
		if err != nil {
			app.errorLog.Printf("signUp: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			app.errorLog.Printf("signUp: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
			return
		}

		err = internal.Registration(app.database, user)
		if err != nil {
			app.errorLog.Printf("signUp: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		app.notice.Exist = false

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) signIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sign-in" {
		app.errorLog.Printf("signIn: invalid path %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("./ui/template/signIn.html")
		if err != nil {
			app.errorLog.Printf("signIn: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, app.notice)
		if err != nil {
			app.errorLog.Printf("signIn: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			app.errorLog.Printf("signIn: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}

		nick := r.FormValue("nickname")
		password := r.FormValue("password")

		err = internal.Login(app.database, nick, password)
		if err != nil {
			app.notice = notification{
				Content: "Nickname or password invalid",
				Exist:   true,
			}
			app.infoLog.Printf("signIn: %s\n", err.Error())
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		sessionToken, err := uuid.NewV4()
		if err != nil {
			app.errorLog.Printf("signIn: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		expiresAt := time.Now().Add(12 * time.Hour)

		err = database.CreateSession(app.database, nick, sessionToken.String(), expiresAt)
		if err != nil {
			app.errorLog.Printf("signIn: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

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
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			app.errorLog.Printf("logOut: unauthorized user\n")
			return
		}

		app.errorLog.Printf("logOut: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return

	}

	sessionToken := c.Value
	err = database.DeleteSession(app.database, sessionToken)
	if err != nil {
		app.errorLog.Printf("logOut: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
