package main

import (
	"log"
	"net/http"
)

func main() {
	// http.NewServeMux initialises a new servemux
	// and is used to register handlers for a URL pattern
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// http.ListenAndServe starts a new web server
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
