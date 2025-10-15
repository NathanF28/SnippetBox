package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"the_Elir.net/internal/models"
)

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
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
	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	// Create a template cache and add to the application struct
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error("unable to create template cache", "err", err)
		os.Exit(1)
	}
	app.templateCache = templateCache

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
