package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the SQLite database
func InitDB(dbPath string) (*sql.DB, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS snapshots (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL,
		timestamp TEXT NOT NULL,
		filename TEXT NOT NULL,
		data TEXT NOT NULL,
		uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(ip, timestamp)
	);

	CREATE INDEX IF NOT EXISTS idx_snapshots_ip ON snapshots(ip);
	CREATE INDEX IF NOT EXISTS idx_snapshots_timestamp ON snapshots(timestamp);
	`

	_, err := db.Exec(schema)
	return err
}
