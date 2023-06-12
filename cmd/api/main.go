package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type config struct {
	bindAddr string
	env      string
}

type application struct {
	config config
	logger *log.Logger
	db     *sql.DB
}

func main() {
	var cfg config
	flag.StringVar(&cfg.bindAddr, "bind", "localhost:4000", "Address to bind to")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB("/home/danimal/projects/te-go/database.sqlite3")
	if err != nil {
		logger.Fatal(err)
	}
	app := &application{
		config: cfg,
		logger: logger,
		db:     db,
	}
	srv := &http.Server{
		Addr:         cfg.bindAddr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
