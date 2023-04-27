package main

import (
	"flag"
	"log"
	"net/http"
)

type Config struct {
	Addr      string
	StaticDir string
}

func main() {
	cfg := new(Config)

	// Define a new command-line flag with the name 'addr', a default value of ":4000" and some short
	// help text explaining what the flag controls
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "../../ui/static/", "Path to static assets")

	// Parse flags before you use them
	flag.Parse()

	// http.NewServeMux initialises a new servemux
	// and is used to register handlers for a URL pattern
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Use the relative path to create a file server at that path
	fs := http.FileServer(http.Dir(cfg.StaticDir))
	// When this handler receives a request, it will remove the leading slash from the URL path and
	// then search the file server directory for the corresponding file to send the user.
	// This is done by stripping the leading "/static" from the URL path before passing it to the file
	// server, otherwise it will be looking for a file which does not exist
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	// http.ListenAndServe starts a new web server
	log.Printf("Starting server on %s\n", cfg.Addr)
	err := http.ListenAndServe(cfg.Addr, mux)
	log.Fatal(err)
}
