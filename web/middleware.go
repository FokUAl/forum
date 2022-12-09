package web

import (
	"forumAA/database"
	"net/http"
	"time"
)

func isExpired(expiry time.Time) bool {
	return expiry.Before(time.Now())
}

func (app *application) checkUser(w http.ResponseWriter, r *http.Request) database.User {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return database.User{}
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return database.User{}

	}
	sessionToken := c.Value

	// We then get the session from our session map
	result, err := database.GetUserByToken(app.database, sessionToken)
	if err != nil {
		// If the session token is not present in sessions table, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return database.User{}
	}

	expiry, err := database.GetExpiryByToken(app.database, sessionToken)
	if err != nil {
		// If the session token is not present in session table, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return database.User{}
	}

	if isExpired(expiry) {
		database.DeleteSession(app.database, sessionToken)
		http.Redirect(w, r, "/", 300)
	}

	return result
}
