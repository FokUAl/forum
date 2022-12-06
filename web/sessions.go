package web

import (
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

// func (app *application) home(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	}

// 	if r.URL.Path != "/" {
// 		http.NotFound(w, r)
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}

// 	c, err := r.Cookie("session_token")
// 	if err != nil {
// 		if err == http.ErrNoCookie {
// 			// If the cookie is not set, return an unauthorized status
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}
// 		// For any other type of error, return a bad request status
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	sessionToken := c.Value

// 	// We then get the session from our session map
// 	userSession, exists := sessions[sessionToken]
// 	if !exists {
// 		// If the session token is not present in session map, return an unauthorized error
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	if userSession.isExpired() {
// 		delete(sessions, sessionToken)

// 		t, err := template.ParseFiles("./ui/template/home.html")
// 		if err != nil {
// 			log.Println(err.Error())
// 			http.Error(w, "File not found: index.html", 500)
// 			return
// 		}

// 		err = t.Execute(w, nil)
// 		if err != nil {
// 			log.Println(err.Error())
// 			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		}
// 	}
// }
