package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os" // New import
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"    // New import
	"github.com/go-playground/form/v4" // New import
	_ "github.com/go-sql-driver/mysql"
	"the_Elir.net/internal/models"
)

type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	//accept network port address and data source name
	addr := flag.String("addr", ":8080", "Network Port Address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQl data source name")
	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := OpenDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	formDecoder := form.NewDecoder()
	sessionManager := scs.New()               // new instance
	sessionManager.Store = mysqlstore.New(db) // where its stored
	sessionManager.Lifetime = 12 * time.Hour  // lifetime

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
	// Create a template cache and add to the application struct
	if err != nil {
		logger.Error("unable to create template cache", "err", err)
		os.Exit(1)
	}

	logger.Info("Starting server", "addr", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func OpenDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close() // close db as not to waste resource if ping fails and connection is open
		return nil, err
	}
	return db, nil
}
