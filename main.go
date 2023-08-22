package main

import (
	"bytes"
	"github.com/freshman-tech/news-demo-starter-files/news"
	"github.com/joho/godotenv"
	"html/template" // html/template pkg -> safe against code injection
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Search represent each query make by the user
type Search struct {
	Query     string
	NextPage  int
	TotalPage int
	Results   *news.Results
}

// parse the index file
var tpl = template.Must(template.ParseFiles("index.html"))

// indexHandler is a simple http request handler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// simply write this into the <body> section
	//w.Write([]byte("<h1>Hello World!</h1>"))

	// update version -> use the template to execute
	// now return the index.html file
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

// searchHandler is for handling search request
func searchHandler(newsapi *news.Client) http.HandlerFunc {
	// return anonymous func
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		// do request
		results, err := newsapi.FetchEverything(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// construct resp
		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		search := &Search{
			Query:     searchQuery,
			NextPage:  nextPage,
			TotalPage: int(math.Ceil(float64(results.TotalResults) / float64(newsapi.PageSize))),
			Results:   results,
		}

		// write response -> back to index.html
		// first write to an empty buffer -> then buffer is written to the ResponseWriter
		// then execute tml directly on ResponseWriter
		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buf.WriteTo(w)
	}
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

	// news api key & set client
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}
	client := &http.Client{Timeout: 10 * time.Second}
	newsapi := news.NewClient(client, apiKey, 20)

	// init a file server by passing the directory where all static files are placed
	fs := http.FileServer(http.Dir("assets"))

	// create new multiplexer, which allowing you to associate different handlers with different URL path
	// which is essentially a router that help you direct incoming HTTP requests to appropriate handler function
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)                        // 'pair' the index path to the handler
	mux.HandleFunc("/search", searchHandler(newsapi))        // 'pair' the search path to the handler
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs)) // pass assets to file system

	http.ListenAndServe(":"+port, mux)

}
