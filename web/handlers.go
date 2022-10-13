package web

import (
	"net/http"
	"html/template"
	"log"
)

func (app *application) home(w http.ResponseWriter, r *http.Request){
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

	t, err := template.ParseFiles("./ui/template/index.html")
	if err != nil{
		log.Println(err.Error())
		http.Error(w, "File not found: index.html", 500)
		return
	}

	err = t.Execute(w, nil)
	if err != nil{
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
