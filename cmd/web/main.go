package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql" // we need the driver's init() function to run so it can register itself with the sql package
	"github.com/mhrdini/snippetbox/internal/models"
)

type Config struct {
	Addr      string
	StaticDir string
	DSN       string
}

// Define an application struct to hold the application-wide dependencies for the web application.
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	cfg := new(Config)

	// Define a new command-line flag with its identifier, a default value, and some short
	// help text explaining what the flag controls
	flag.StringVar(&cfg.Addr, "addr", ":8000", "HTTP network address")
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
	defer db.Close()

	// Initialise template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialise form decoder
	formDecoder := form.NewDecoder()

	// Initialise session manager, configured to use MySQL DB as session store
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Initialise a new instance of application containing the dependencies.
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initialise TLS config for non-default TLS/HTTPS settings
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}, // elliptic curves restriction
		// -- min or max TLS version
		// MinVersion: tls.VersionTLS512,
		// MaxVersion: tls.VersionTLS512,
		// -- cipher suites that aren't weak or those that use ECDHE (forward secrecy)
		// CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384},
	}

	// Initialise a new http.Server struct, setting the Addr and Handler fields to have it use the
	// appropriate network address and routes, and the ErrorLog field so that the server now uses
	// the custom errorLog logger in the event of any problems
	srv := &http.Server{
		Addr:         cfg.Addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(cfg), // servemux in routes.go
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,     // reduce keep-alive to close connections earlier
		ReadTimeout:  5 * time.Second, // close connection if accepted connection still hasn't read req headers/body
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s\n", cfg.Addr)
	err = srv.ListenAndServeTLS("../../tls/cert.pem", "../../tls/key.pem")
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
