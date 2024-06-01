package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/mhrdini/snippetbox/internal/models"
	"github.com/mhrdini/snippetbox/ui"
)

// Define a templateData type to act as the holding structure for any dynamic data
// that we want to pass to our HTML templates.
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any // to pass the validation errors and previously submitted data back to the template when we redisplay the form
	Toast           string
	IsAuthenticated bool
}

func prettyDate(t time.Time) string {
	// Return the empty string id time has the zero value.
	if t.IsZero() {
		return ""
	}

	// Convert the time to UTC before formatting it.
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"prettyDate": prettyDate,
}

// Caches all the templates in a map by using filepath.Glob() to get a slice of all pages
// and then using *template.Template.ParseFiles() and *template.Template.ParseGlob() methods
// to add templates to the template set.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with the extension '.page.html'.
	// This essentially gives us a slice of all the 'page' templates for the application.
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		filename := filepath.Base(page)

		files := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		ts, err := template.New(filename).Funcs(functions).ParseFS(ui.Files, files...)
		if err != nil {
			return nil, err
		}

		cache[filename] = ts
	}

	return cache, nil
}
