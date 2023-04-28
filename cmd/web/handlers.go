package main

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/mhrdini/snippetbox/pkg/models"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body
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

	// Use template.ParseFiles()
	files := []string{
		"../../ui/html/home.page.html",
		"../../ui/html/base.layout.html",
		"../../ui/html/footer.partial.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	// Use Execute() method on the template set to write the template content as the response body.
	// The last parameter to Execute() is the dynamic data that we may pass in.
	err = ts.Execute(w, data)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}

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

	files := []string{
		"../../ui/html/show.page.html",
		"../../ui/html/base.layout.html",
		"../../ui/html/footer.partial.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
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
