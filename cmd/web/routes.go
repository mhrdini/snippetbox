package main

import "net/http"

func (app *application) routes(cfg *Config) *http.ServeMux {
	// http.NewServeMux initialises a new servemux
	// and is used to register handlers for a URL pattern
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Use the relative path to create a file server at that path
	fs := http.FileServer(http.Dir(cfg.StaticDir))
	// When this handler receives a request, it will remove the leading slash from the URL path and
	// then search the file server directory for the corresponding file to send the user.
	// This is done by stripping the leading "/static" from the URL path before passing it to the file
	// server, otherwise it will be looking for a file which does not exist
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	return mux
}
