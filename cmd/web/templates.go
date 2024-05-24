package main

import (
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/mhrdini/snippetbox/internal/models"
	"github.com/mhrdini/snippetbox/ui"
)

// Define a templateData type to act as the holding structure for any dynamic data
// that we want to pass to our HTML templates.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
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

		ts, err := template.New(filename).ParseFS(ui.Files, files...)
		if err != nil {
			return nil, err
		}

		cache[filename] = ts
	}

	return cache, nil
}
