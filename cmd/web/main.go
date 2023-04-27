package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Addr      string
	StaticDir string
}

// Define an application struct to hold the application-wide dependencies for the web application.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	cfg := new(Config)

	// Define a new command-line flag with the name 'addr', a default value of ":4000" and some short
	// help text explaining what the flag controls
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "../../ui/static/", "Path to static assets")

	// Parse flags before you use them
	flag.Parse()

	// Create loggers for writing info and error messages.
	// InfoLog -> Writes to os.Stdout, uses INFO prefix, and flags to include additional info such as
	// local datetime
	// ErrorLog -> Writes to os.Stderr, use log.Lshortfile flag to include relevant file name and line number
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialise a new instance of application containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

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

	// Initialise a new http.Server struct, setting the Addr and Handler fields to have it use the
	// same network address and routes as before, and the ErrorLog field so that the server now uses
	// the custom errorLog logger in the event of any problems
	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	infoLog.Printf("Starting server on %s\n", cfg.Addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
