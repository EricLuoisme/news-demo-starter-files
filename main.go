package main

import (
	"github.com/joho/godotenv"
	"html/template" // html/template pkg -> safe against code injection
	"log"
	"net/http"
	"os"
)

// parse the index file
var tpl = template.Must(template.ParseFiles("index.html"))

// indexHandler is a simple http request handler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// simply write this into the <body> section
	//w.Write([]byte("<h1>Hello World!</h1>"))

	// update version -> use the template to execute
	// now return the index.html file
	tpl.Execute(w, nil)
}

func main() {

	// use .env reading library -> get config from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file!")
	}

	// setting port
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	// init a file server by passing the directory where all static files are placed
	fs := http.FileServer(http.Dir("assets"))

	// create new multiplexer, which allowing you to associate different handlers with different URL path
	// which is essentially a router that help you direct incoming HTTP requests to appropriate handler function
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)                        // 'pair' the path to the handler
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs)) // pass assets to file system

	http.ListenAndServe(":"+port, mux)

}
