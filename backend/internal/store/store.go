package store

import (
	"backend/internal/config"
	"log"
)

// Store is kept for backward compatibility
// It's recommended to use DBStore directly for new code
type Store struct {
	// Embed DBStore to inherit all its methods
	*DBStore
}

// New creates a new Store for backward compatibility
// It uses the new DBStore implementation with PostgreSQL
func New(dataSourceName string) (*Store, error) {
	log.Printf("Store.New: Creating a new store with PostgreSQL using path: %s", dataSourceName)

	// Create a default config that uses PostgreSQL
	cfg := &config.DatabaseConfig{
		Host:     "postgres",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "chatapp",
		SSLMode:  "disable",
	}

	// Create a new DBStore with the config
	dbStore, err := NewDBStore(cfg)
	if err != nil {
		return nil, err
	}

	// Wrap the DBStore in a Store for backward compatibility
	return &Store{DBStore: dbStore}, nil
}

// Migrate delegates to the embedded DBStore's Migrate method
func (s *Store) Migrate() error {
	log.Println("Store.Migrate: Delegating to DBStore.Migrate")
	return s.DBStore.Migrate()
}
