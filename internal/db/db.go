package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
}

func InitDatabase(dbPath string) (*Database, error) {
	_, err := os.Stat(dbPath)
	dbExists := !errors.Is(err, os.ErrNotExist)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	database := &Database{db}
	if !dbExists {
		if err := database.initSchema(); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to initialize schema: %w", err)
		}
	}
	return database, nil
}

func (db *Database) initSchema() error {
	// Create users table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Add any additional schema initialization here
	// For example, creating indexes:
	createIndexSQL := `
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	`

	_, err = db.Exec(createIndexSQL)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

func (db *Database) CheckHealth() error {
	return db.Ping()
}
