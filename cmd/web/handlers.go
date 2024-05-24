package main

import (
	"net/http"
	"strconv"

	"github.com/mhrdini/snippetbox/internal/models"
)

// Define a home handler function
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	// for _, snippet := range s {
	// 	fmt.Fprintf(w, "%v\n", snippet)
	// }

	data := &templateData{Snippets: s}

	app.render(w, r, "home.tmpl.html", data)
}

// Add a showSnippet handler function that receives an id query parameter
// that must be an integer greater than or equal to 1.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippet: s}
	app.render(w, r, "show.tmpl.html", data)
}

// Add a createSnippet handler function that only receives POST requests.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// id, err := app.snippets.Insert(title, content, expires)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
	// http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

	w.Write([]byte("Create a new snippet..."))
}
