package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func NewConnection(dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

func RunMigrations(db *sqlx.DB, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db.DB, dir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
