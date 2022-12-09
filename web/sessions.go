package web

import (
	"forumAA/database"
	"net/http"
	"time"
)

var sessions = map[string]session{}

type session struct {
	username string
	expiry   time.Time
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func (app *application) checkUser(w http.ResponseWriter, r *http.Request) database.User {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
		} else {
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return database.User{}
		}
	}
	sessionToken := c.Value

	// We then get the session from our session map
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return database.User{}
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
		http.Redirect(w, r, "localhost:4888/signin", 300)
	}

	return database.User{}
}
