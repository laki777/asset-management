package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var tpl *template.Template

type pageData struct {
	Title      string
	FirstName  string
	SecondName string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {

	gorilaMux := mux.NewRouter()

	aboutHandler := http.HandlerFunc(about)
	indexHandler := http.HandlerFunc(index)
	signupHandler := http.HandlerFunc(signup)
	searchHandler := http.HandlerFunc(search)
	logoutHandler := http.HandlerFunc(logout)
	contactHandler := http.HandlerFunc(contact)
	insertHandler := http.HandlerFunc(insert)

	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	gorilaMux.Handle("/", handlers.LoggingHandler(logFile, handlers.CompressHandler(indexHandler)))
	gorilaMux.Handle("/about", handlers.LoggingHandler(logFile, handlers.CompressHandler(aboutHandler)))
	gorilaMux.Handle("/signup", handlers.LoggingHandler(logFile, handlers.CompressHandler(signupHandler)))
	gorilaMux.Handle("/search", handlers.LoggingHandler(logFile, handlers.CompressHandler(searchHandler)))
	gorilaMux.Handle("/logout", handlers.LoggingHandler(logFile, handlers.CompressHandler(logoutHandler)))
	gorilaMux.Handle("/contact", handlers.LoggingHandler(logFile, handlers.CompressHandler(contactHandler)))
	gorilaMux.Handle("/insert", handlers.LoggingHandler(logFile, handlers.CompressHandler(insertHandler)))

	http.Handle("/", gorilaMux)

	http.Handle("/imgs/", http.StripPrefix("/imgs", http.FileServer(http.Dir("./imgs"))))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)

}
