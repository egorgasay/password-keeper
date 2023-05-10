package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"password-keeper/internal/storage/sqllike"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// File driver for migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// Postgres driver
	_ "github.com/lib/pq"
)

// Postgres struct with *sql.DB instance.
// It has methods for working with URLs.
type Postgres struct {
	sqllike.DB
}

// New Postgres struct constructor.
func New(db *sql.DB, path string) (*Postgres, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("can't init migrate instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		path,
		"postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("can't create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("can't migrate up: %w", err)
	}

	return &Postgres{DB: sqllike.DB{DB: db}}, nil
}
