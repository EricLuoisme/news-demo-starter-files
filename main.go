package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// simply write this into the <body> section
	w.Write([]byte("<h1>Hello World!</h1>"))
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

	// create new multiplexer, which allowing you to associate different handlers with different URL path
	// which is essentially a router that help you direct incoming HTTP requests to appropriate handler function
	mux := http.NewServeMux()

	// 'pair' the path to the handler
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)

}
