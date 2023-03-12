package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Vesino/todos/internal/data"
	"github.com/Vesino/todos/internal/jsonlog"
	"github.com/Vesino/todos/internal/mailer"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/wait"
)

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
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func TestApplication(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	const dbname = "test-db"
	const user = "postgres"
	const password = "password"

	port, err := nat.NewPort("tcp", "5432")
	require.NoError(t, err)

	container, err := startContainer(ctx,
		WithPort(port.Port()),
		WithInitialDatabase(user, password, dbname),
		WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	containerPort, err := container.MappedPort(ctx, port)
	assert.NoError(t, err)

	host, err := container.Host(ctx)
	assert.NoError(t, err)

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, containerPort.Port(), user, password, dbname)

	// perform assertions
	// db, err := sql.Open("postgres", connStr)

	cfg := config{
		env: "test",
		db: struct{ dsn string }{
			dsn: connStr,
		},
		smtp: struct {
			host     string
			port     int
			username string
			password string
			sender   string
		}{
			host:     "smtp.mailtrap.io",
			port:     2525,
			username: "your-username",
			password: "your-password",
			sender:   "Testing <test@example.com>",
		},
	}

	db, err := openDB(cfg)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, db)
	defer db.Close()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	if app == nil {
		t.Fatal("App is nil")
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
