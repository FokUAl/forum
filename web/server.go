package web

import (
	"log"
	"net/http"
	"os"
)

type application struct {
	server   *http.Server
	errorLog *log.Logger
	infoLog  *log.Logger
}

func Run(port string) {
	mux := http.NewServeMux()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	app := &application{
		server:   srv,
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// TO DO
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/post", app.post)
	mux.HandleFunc("/signup", app.signUp)

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	infoLog.Printf("Launch server on %s", port)

	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
