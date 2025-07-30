package store

import (
	"backend/internal/models"
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // The SQLite driver
)

// Store manages the database connection and operations.
type Store struct {
	db *sql.DB
}

// New creates a new Store and initializes the database connection.
func New(dataSourceName string) (*Store, error) {
	log.Printf("Store.New: Connecting to database: %s", dataSourceName)
	
	// Make sure the database file path is valid
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Printf("Store.New: Failed to open database: %v", err)
		return nil, err
	}
	
	// Set connection parameters for better reliability
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)
	
	// Verify connection is working
	if err := db.Ping(); err != nil {
		log.Printf("Store.New: Failed to ping database: %v", err)
		return nil, err
	}
	
	log.Println("Store.New: Database connection established successfully")
	return &Store{db: db}, nil
}

// Migrate creates the necessary database tables if they don't already exist.
func (s *Store) Migrate() error {
	log.Println("Store.Migrate: Starting database migration")
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
	if err != nil {
		log.Printf("Store.Migrate: Database migration failed: %v", err)
		return err
	}
	
	// Verify tables were created by checking if we can query them
	log.Println("Store.Migrate: Verifying tables were created")
	tables := []string{"users", "chat_rooms", "messages"}
	for _, table := range tables {
		query := "SELECT name FROM sqlite_master WHERE type='table' AND name=?"
		row := s.db.QueryRow(query, table)
		var name string
		err := row.Scan(&name)
		if err != nil {
			log.Printf("Store.Migrate: Table verification failed for %s: %v", table, err)
		} else {
			log.Printf("Store.Migrate: Table %s verified", name)
		}
	}
	
	log.Println("Store.Migrate: Database migration completed successfully")
	return nil
}

// CreateUser inserts a new user into the database.
func (s *Store) CreateUser(user *models.User) error {
	log.Printf("Store.CreateUser: Inserting user with ID: %s, Username: %s", user.ID, user.Username)
	const query = `INSERT INTO users (id, username, password) VALUES (?, ?, ?)`
	_, err := s.db.Exec(query, user.ID, user.Username, user.Password)
	if err != nil {
		log.Printf("Store.CreateUser: Database error: %v", err)
		return err
	}
	log.Printf("Store.CreateUser: User successfully inserted into database")
	return nil
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

// SaveMessage saves a message to the database.
func (s *Store) SaveMessage(message *models.Message) error {
	const query = `INSERT INTO messages (id, room_id, sender_id, content, timestamp) VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, message.ID, message.RoomID, message.SenderID, message.Content, message.Timestamp)
	return err
}

// GetMessagesByRoom retrieves messages for a specific room with pagination.
func (s *Store) GetMessagesByRoom(roomID string, limit, offset int) ([]models.Message, error) {
	const query = `
		SELECT m.id, m.room_id, m.sender_id, m.content, m.timestamp, u.username
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.room_id = ?
		ORDER BY m.timestamp DESC
		LIMIT ? OFFSET ?
	`
	rows, err := s.db.Query(query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var username string
		if err := rows.Scan(&msg.ID, &msg.RoomID, &msg.SenderID, &msg.Content, &msg.Timestamp, &username); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// GetMessagesSince retrieves messages for a room since a specific time.
func (s *Store) GetMessagesSince(roomID string, since time.Time) ([]models.Message, error) {
	const query = `
		SELECT m.id, m.room_id, m.sender_id, m.content, m.timestamp, u.username
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.room_id = ? AND m.timestamp > ?
		ORDER BY m.timestamp ASC
	`
	rows, err := s.db.Query(query, roomID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var username string
		if err := rows.Scan(&msg.ID, &msg.RoomID, &msg.SenderID, &msg.Content, &msg.Timestamp, &username); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
