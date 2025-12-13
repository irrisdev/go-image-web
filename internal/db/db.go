package db

import (
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	sqlite3m "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	MigrationsDir string = "file://internal/db/migrations"
)

func EnsureSchema(db *sqlx.DB) error {
	driver, err := sqlite3m.WithInstance(db.DB, &sqlite3m.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(MigrationsDir, "sqlite3", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
