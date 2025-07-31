package config

import (
	"fmt"
	"os"
)

// DatabaseConfig holds the database connection parameters
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDatabaseConfig creates a new database configuration from environment variables
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "chatapp"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// ConnectionString returns the PostgreSQL connection string
func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// SQLiteConnectionString returns the SQLite connection string for local development
func (c *DatabaseConfig) SQLiteConnectionString() string {
	return "./chatapp.db"
}

// IsSQLite returns true if the database is SQLite
func (c *DatabaseConfig) IsSQLite() bool {
	// Check if DB_TYPE environment variable is set to sqlite
	return getEnv("DB_TYPE", "sqlite") == "sqlite"
}

// Helper function to get environment variables with default values
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
