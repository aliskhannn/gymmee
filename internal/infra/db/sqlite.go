package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Import sqlite3 driver
)

// NewSqliteDB creates a new database connection and verifies it.
// It also ensures that foreign key constraints are enabled for SQLite.
func NewSqliteDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sqlite: %w", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}
