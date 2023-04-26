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

	// Use the relative path to create a file server at that path
	fs := http.FileServer(http.Dir("../../ui/static/"))
	// When this handler receives a request, it will remove the leading slash from the URL path and
	// then search the file server directory for the corresponding file to send the user.
	// This is done by stripping the leading "/static" from the URL path before passing it to the file
	// server, otherwise it will be looking for a file which does not exist
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	// http.ListenAndServe starts a new web server
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
