package store

import (
	"backend/internal/models"
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // The SQLite driver
)

// Store manages the database connection and operations.
type Store struct {
	db *sql.DB
}

// New creates a new Store and initializes the database connection.
func New(dataSourceName string) (*Store, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// Migrate creates the necessary database tables if they don't already exist.
func (s *Store) Migrate() error {
	const schema = `
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
	_, err := s.db.Exec(schema)
	return err
}

// CreateUser inserts a new user into the database.
func (s *Store) CreateUser(user *models.User) error {
	const query = `INSERT INTO users (id, username, password) VALUES (?, ?, ?)`
	_, err := s.db.Exec(query, user.ID, user.Username, user.Password)
	return err
}

// GetUserByUsername retrieves a user from the database by username.
func (s *Store) GetUserByUsername(username string) (*models.User, error) {
	const query = `SELECT id, username, password FROM users WHERE username = ?`
	row := s.db.QueryRow(query, username)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateRoom inserts a new chat room into the database.
func (s *Store) CreateRoom(room *models.ChatRoom) error {
	const query = `INSERT INTO chat_rooms (id, name) VALUES (?, ?)`
	_, err := s.db.Exec(query, room.ID, room.Name)
	return err
}

// GetAllRooms retrieves all chat rooms from the database.
func (s *Store) GetAllRooms() ([]models.ChatRoom, error) {
	const query = `SELECT id, name FROM chat_rooms`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.ChatRoom
	for rows.Next() {
		var room models.ChatRoom
		if err := rows.Scan(&room.ID, &room.Name); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}
