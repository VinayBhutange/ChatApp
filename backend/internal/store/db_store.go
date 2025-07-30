package store

import (
	"backend/internal/config"
	"database/sql"
	"fmt"
	"log"
	"time"

	// Import database drivers
	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver (kept for backward compatibility)
)

// DBStore implements the StoreInterface with a SQL database
type DBStore struct {
	db     *sql.DB
	config *config.DatabaseConfig
}

// Ensure DBStore implements StoreInterface
var _ StoreInterface = (*DBStore)(nil)

// NewDBStore creates a new database store
func NewDBStore(cfg *config.DatabaseConfig) (*DBStore, error) {
	var db *sql.DB
	var err error

	// Use PostgreSQL by default
	log.Println("Using PostgreSQL database")
	db, err = sql.Open("postgres", cfg.ConnectionString())

	// Keep SQLite as a fallback option if explicitly configured
	if cfg.IsSQLite() {
		log.Println("Using SQLite database")
		db, err = sql.Open("sqlite3", cfg.SQLiteConnectionString())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection parameters for better reliability
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Verify connection is working
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return &DBStore{db: db, config: cfg}, nil
}

// Close closes the database connection
func (s *DBStore) Close() error {
	return s.db.Close()
}

// Migrate creates the necessary database tables if they don't already exist
func (s *DBStore) Migrate() error {
	log.Println("Starting database migration")

	// Always use PostgreSQL schema
	schema := postgresMigrationSchema

	_, err := s.db.Exec(schema)
	if err != nil {
		log.Printf("Database migration failed: %v", err)
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}

// SQLite migration schema
const sqliteMigrationSchema = `
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS chat_rooms (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    room_id TEXT NOT NULL,
    sender_id TEXT NOT NULL,
    content TEXT NOT NULL,
    timestamp DATETIME NOT NULL,
    FOREIGN KEY (room_id) REFERENCES chat_rooms(id),
    FOREIGN KEY (sender_id) REFERENCES users(id)
);
`

// PostgreSQL migration schema
const postgresMigrationSchema = `
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS chat_rooms (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    room_id TEXT NOT NULL,
    sender_id TEXT NOT NULL,
    content TEXT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (room_id) REFERENCES chat_rooms(id),
    FOREIGN KEY (sender_id) REFERENCES users(id)
);
`
