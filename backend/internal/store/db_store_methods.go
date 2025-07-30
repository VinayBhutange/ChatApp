package store

import (
	"backend/internal/models"
	"time"
)

// CreateUser inserts a new user into the database.
func (s *DBStore) CreateUser(user *models.User) error {
	var query string
	if s.config.IsSQLite() {
		query = `INSERT INTO users (id, username, password) VALUES (?, ?, ?)`
	} else {
		query = `INSERT INTO users (id, username, password) VALUES ($1, $2, $3)`
	}
	
	_, err := s.db.Exec(query, user.ID, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByUsername retrieves a user from the database by username.
func (s *DBStore) GetUserByUsername(username string) (*models.User, error) {
	var query string
	if s.config.IsSQLite() {
		query = `SELECT id, username, password FROM users WHERE username = ?`
	} else {
		query = `SELECT id, username, password FROM users WHERE username = $1`
	}
	
	row := s.db.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateRoom creates a new chat room in the database.
func (s *DBStore) CreateRoom(room *models.ChatRoom) error {
	var query string
	if s.config.IsSQLite() {
		query = `INSERT INTO chat_rooms (id, name) VALUES (?, ?)`
	} else {
		query = `INSERT INTO chat_rooms (id, name) VALUES ($1, $2)`
	}
	
	_, err := s.db.Exec(query, room.ID, room.Name)
	if err != nil {
		return err
	}
	return nil
}

// GetAllRooms retrieves all chat rooms from the database.
func (s *DBStore) GetAllRooms() ([]*models.ChatRoom, error) {
	// This query is the same for both SQLite and PostgreSQL
	const query = `SELECT id, name FROM chat_rooms`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*models.ChatRoom
	for rows.Next() {
		var room models.ChatRoom
		if err := rows.Scan(&room.ID, &room.Name); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

// SaveMessage saves a message to the database.
func (s *DBStore) SaveMessage(message *models.Message) error {
	var query string
	if s.config.IsSQLite() {
		query = `INSERT INTO messages (id, room_id, sender_id, content, timestamp) VALUES (?, ?, ?, ?, ?)`
	} else {
		query = `INSERT INTO messages (id, room_id, sender_id, content, timestamp) VALUES ($1, $2, $3, $4, $5)`
	}
	
	_, err := s.db.Exec(query, message.ID, message.RoomID, message.SenderID, message.Content, message.Timestamp)
	if err != nil {
		return err
	}
	return nil
}

// GetMessagesByRoom retrieves messages for a specific room from the database.
func (s *DBStore) GetMessagesByRoom(roomID string) ([]*models.Message, error) {
	var query string
	if s.config.IsSQLite() {
		query = `SELECT id, room_id, sender_id, content, timestamp FROM messages WHERE room_id = ? ORDER BY timestamp DESC LIMIT 50`
	} else {
		query = `SELECT id, room_id, sender_id, content, timestamp FROM messages WHERE room_id = $1 ORDER BY timestamp DESC LIMIT 50`
	}
	
	rows, err := s.db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.RoomID, &msg.SenderID, &msg.Content, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// GetMessagesSince retrieves messages for a specific room since a given time.
func (s *DBStore) GetMessagesSince(roomID string, since time.Time) ([]*models.Message, error) {
	var query string
	if s.config.IsSQLite() {
		query = `SELECT id, room_id, sender_id, content, timestamp FROM messages WHERE room_id = ? AND timestamp > ? ORDER BY timestamp ASC`
	} else {
		query = `SELECT id, room_id, sender_id, content, timestamp FROM messages WHERE room_id = $1 AND timestamp > $2 ORDER BY timestamp ASC`
	}
	
	rows, err := s.db.Query(query, roomID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.RoomID, &msg.SenderID, &msg.Content, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
