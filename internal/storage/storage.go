package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func ConnectDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	migrationsDir := "./migrations"

	if err = initMigrations(db, migrationsDir); err != nil {
		return nil, err
	}

	return db, nil
}

// Initializing database migration
func initMigrations(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(db, dir)
	if err != nil {
		return err
	}

	return nil
}
