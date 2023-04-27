package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // we need the driver's init() function to run so it can register itself with the sql package
)

type Config struct {
	Addr      string
	StaticDir string
	DSN       string
}

// Define an application struct to hold the application-wide dependencies for the web application.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	cfg := new(Config)

	// Define a new command-line flag with its identifier, a default value and some short
	// help text explaining what the flag controls
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "../../ui/static/", "Path to static assets")
	flag.StringVar(&cfg.DSN, "dsn", "web:web@/snippetbox?parseTime=true", "MySQL database connection string")

	// Parse flags before you use them
	flag.Parse()

	// Create loggers for writing info and error messages.
	// InfoLog -> Writes to os.Stdout, uses INFO prefix, and flags to include additional info such as
	// local datetime
	// ErrorLog -> Writes to os.Stderr, use log.Lshortfile flag to include relevant file name and line number
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Open DB here
	db, err := openDB(cfg.DSN)
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialise a new instance of application containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Initialise a new http.Server struct, setting the Addr and Handler fields to have it use the
	// appropriate network address and routes, and the ErrorLog field so that the server now uses
	// the custom errorLog logger in the event of any problems
	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  app.routes(cfg), // servemux in routes.go
	}
	infoLog.Printf("Starting server on %s\n", cfg.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// openDB() wraps sql.Open() and returns a sql.DB connection pool for a given DSN
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
