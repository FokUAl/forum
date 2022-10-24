package web

import (
	"html/template"
	"log"
	"net/http"
)

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
		http.Error(w, "File not found: index.html", 500)
		return
	}

	err = t.Execute(w, nil)
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

	if r.URL.Path != "/post" {
		http.NotFound(w, r)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles("./ui/template/post.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "File not found: post.html", 500)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/signup" {
		http.NotFound(w, r)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles("./ui/template/signup.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "File not found: signup.html", 500)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
