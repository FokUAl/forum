package web

import (
	"database/sql"
	"forumAA/database"
	"log"
	"net/http"
	"os"
)

type application struct {
	server   *http.Server
	errorLog *log.Logger
	infoLog  *log.Logger
	database *sql.DB
}

func Run(port string) error {
	mux := http.NewServeMux()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	dbHandler := database.Init()

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	app := &application{
		server:   srv,
		errorLog: errorLog,
		infoLog:  infoLog,
		database: dbHandler,
	}

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/post/", app.post)
	mux.HandleFunc("/sign-up", app.signUp)
	mux.HandleFunc("/sign-in", app.signIn)
	mux.HandleFunc("/profile/", app.profile)
	mux.HandleFunc("/logout", app.logOut)
	mux.HandleFunc("/create-post", app.createPost)
	mux.HandleFunc("/comment/like/", app.likeComment)
	mux.HandleFunc("/post/like/", app.likePost)

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	infoLog.Printf("Launch server on %s", port)

	err := srv.ListenAndServe()
	errorLog.Fatal(err)

	return err
}
