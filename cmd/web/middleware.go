package main

import (
	"fmt"
	"net/http"
)

// Middleware can be executed:
// - before servemux -> applied to all requests
// - after servemux -> applied to some requests

// This is set either in routes or main

// /* Middleware template: */
// func myMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		/* any code here will execute on the way down the chain */
//    /* any early returns here will revert control back upstream, this pattern is used to restrict access (e.g. authorization) */
// 		next.ServeHTTP(w, r)
// 		/* any code here will execute on the way up the chain */
// 	})
// }

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

// Use recover() to check if a panic happened and defer to call serverError when it has panicked
// NOTE: Will only recover panics happening in the same goroutine that executed this middleware
// Need to use the following pattern to recover from panics within a defer func within a goroutine
// func: if err := recover(); err != nil { log.Print(fmt.Errorf("%s\n%s", err, debug.Stack())) }
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Setting this acts as a trigger to make Go's HTTP/1 server auto-close the connection
				w.Header().Set("Connection", "close")
				// Normalise any-typed error from recover() into an Errorf object format
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
