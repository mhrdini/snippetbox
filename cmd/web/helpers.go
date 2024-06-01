package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
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
func (app *application) render(w http.ResponseWriter, status int, filename string, data *templateData) {

	ts, ok := app.templateCache[filename]
	if !ok {
		app.serverError(w, fmt.Errorf("not found: template %s does not exist", filename))
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

// Create a newTemplateData() helper, which returns a pointer to a templateData struct
// initialized with the current year.
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Toast:           app.sessionManager.PopString(r.Context(), "toast"),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

// Create a decodePostForm() helper method. dst is the target destination
// that we want to decode the form data into
func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call the Decode() method to insert decoded form data into dst from the request form data.
	// i.e. fill the dst  with the relevant values from the HTML form.
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// When target destination is not a non-nil pointer, Decode() method will return
		// an error with the type *form.InvalidDecoderError and we check for this using errors.As()
		// and raise a panic rather than returning the error
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// For all other errors
		return err
	}

	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
