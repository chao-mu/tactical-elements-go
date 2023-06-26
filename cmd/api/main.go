package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chao-mu/tactical-elements-go/internal/data"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

type config struct {
	bindAddr string
	env      string
	db       data.DBConfig
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
	flag.StringVar(&cfg.db.DSN, "db-dsn", os.Getenv("TE_DB_DSN"), "Database DSN")
	flag.IntVar(&cfg.db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := data.OpenDB(cfg.db)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

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
