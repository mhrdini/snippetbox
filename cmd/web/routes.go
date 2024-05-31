package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/mhrdini/snippetbox/ui"
)

func (app *application) routes(cfg *Config) http.Handler {
	// httprouter.New initialises a new servemux
	// and is used to register handlers for a URL pattern
	router := httprouter.New()

	// Createa handler function which wraps our app.notFound helper,
	// then assign it as the custom handler for 404 Not Found responses.
	// Can also be done for 405 Method Not Allowed using router.MethodNotAllowed.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Patterns may include:
	// - :named parameters, as a wildcard
	// - *catch-all parameters, matches everything, should be at the end of a filepath

	// Use the relative path to create a file server at that path
	fs := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fs)

	dynamic := alice.New(app.sessionManager.LoadAndSave)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.viewSnippet))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.createSnippet))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.createSnippetPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
