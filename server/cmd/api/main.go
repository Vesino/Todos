package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Vesino/todos/internal/data"
	"github.com/Vesino/todos/internal/jsonlog"
	"github.com/Vesino/todos/internal/mailer"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
		// maxOpenConns int
		// maxIdleConns int
		// maxIdleTime  string
	}
	smtp struct {
		host string
		port int
		username string
		password string
		sender string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 3333, "Api server port")
	flag.StringVar(&cfg.env, "enviroment", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")

	// Read the SMTP server configuration settings into the config struct, using the
	// settings as the default values.
	flag.StringVar(&cfg.smtp.host, "smtp-host", "0.0.0.0", "SMTP Host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 1025, "SMTP Port")
	flag.StringVar(&cfg.smtp.username, "username", "", "SMTP Username")
	flag.StringVar(&cfg.smtp.password, "password", "", "SMTP Password")
	flag.StringVar(&cfg.smtp.sender, "sender", "Todos <no-reply@greenlight.alexedwards.net>", "SMTP sender")


	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	logger.PrintInfo("database connection pool stablished", nil)

	app := application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()

	if err != nil {
		fmt.Println(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
