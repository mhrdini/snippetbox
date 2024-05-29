package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// Uses debug.Stack() function to get a stack trace for the current goroutine
// and appends it to the log message. Set the frame depth to 2, making it report the file name and
// line number one step back in the stack trace (instead of showing helpers.go because this was
// where the error was called in)
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Uses http.StatusText() to automaticaly generate a human-friendly text representation of a given
// HTTP status code
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Retrieve appropriate template from cache set based on page name, if not found then return server
// error helper method.
func (app *application) render(w http.ResponseWriter, status int, filename string, td *templateData) {
	ts, ok := app.templateCache[filename]
	if !ok {
		app.serverError(w, fmt.Errorf("not found: template %s does not exist", filename))
		return
	}

	err := ts.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
	}
}

// Create a newTemplateData() helper, which returns a pointer to a templateData struct
// initialized with the current year.
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}
