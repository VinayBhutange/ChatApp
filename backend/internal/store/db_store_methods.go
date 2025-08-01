package store

import (
	"backend/internal/models"
	"database/sql"
	"time"
)

// CreateUser creates a new user in the database.
func (s *DBStore) CreateUser(user *models.User) error {
	query := `INSERT INTO users (id, username, password) VALUES (?, ?, ?)`
	if !s.config.IsSQLite() {
		query = `INSERT INTO users (id, username, password) VALUES ($1, $2, $3)`
	}
	_, err := s.db.Exec(query, user.ID, user.Username, user.Password)
	return err
}

// GetUserByUsername retrieves a user by their username.
func (s *DBStore) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password FROM users WHERE username = ?`
	if !s.config.IsSQLite() {
		query = `SELECT id, username, password FROM users WHERE username = $1`
	}

	var user models.User
	err := s.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found is not an error
		}
		return nil, err
	}
	return &user, nil
}

// CreateRoom creates a new chat room in the database.
func (s *DBStore) CreateRoom(room *models.ChatRoom) error {
	query := `INSERT INTO chat_rooms (id, name, owner_id, room_type) VALUES (?, ?, ?, ?)`
	if !s.config.IsSQLite() {
		query = `INSERT INTO chat_rooms (id, name, owner_id, room_type) VALUES ($1, $2, $3, $4)`
	}

	_, err := s.db.Exec(query, room.ID, room.Name, room.OwnerID, room.RoomType)
	return err
}

// AddRoomMember adds a user to a room with a specific status.
func (s *DBStore) AddRoomMember(member *models.RoomMember) error {
	query := `INSERT INTO room_members (room_id, user_id, status) VALUES (?, ?, ?)`
	if !s.config.IsSQLite() {
		query = `INSERT INTO room_members (room_id, user_id, status) VALUES ($1, $2, $3)`
	}

	_, err := s.db.Exec(query, member.RoomID, member.UserID, member.Status)
	return err
}

// GetRoomsByUserID fetches all public rooms and private rooms the user is a member of.
func (s *DBStore) GetRoomsByUserID(userID string) ([]*models.ChatRoom, error) {
	query := `
		SELECT DISTINCT cr.id, cr.name, cr.owner_id, cr.room_type
		FROM chat_rooms cr
		LEFT JOIN room_members rm ON cr.id = rm.room_id
		WHERE cr.room_type = 'public' OR (rm.user_id = ? AND rm.status = 'member')
	`
	if !s.config.IsSQLite() {
		query = `
			SELECT DISTINCT cr.id, cr.name, cr.owner_id, cr.room_type
			FROM chat_rooms cr
			LEFT JOIN room_members rm ON cr.id = rm.room_id
			WHERE cr.room_type = 'public' OR (rm.user_id = $1 AND rm.status = 'member')
		`
	}

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*models.ChatRoom
	for rows.Next() {
		var room models.ChatRoom
		if err := rows.Scan(&room.ID, &room.Name, &room.OwnerID, &room.RoomType); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}

	return rooms, rows.Err()
}

// SaveMessage saves a new message to the database.
func (s *DBStore) SaveMessage(message *models.Message) error {
	query := `INSERT INTO messages (id, room_id, sender_id, content, timestamp) VALUES (?, ?, ?, ?, ?)`
	if !s.config.IsSQLite() {
		query = `INSERT INTO messages (id, room_id, sender_id, content, timestamp) VALUES ($1, $2, $3, $4, $5)`
	}

	_, err := s.db.Exec(query, message.ID, message.RoomID, message.SenderID, message.Content, message.Timestamp)
	return err
}

// GetMessagesByRoom retrieves all messages for a specific room.
func (s *DBStore) GetMessagesByRoom(roomID string) ([]*models.Message, error) {
	query := `SELECT id, room_id, sender_id, content, timestamp FROM messages WHERE room_id = ? ORDER BY timestamp ASC`
	if !s.config.IsSQLite() {
		query = `SELECT id, room_id, sender_id, content, timestamp FROM messages WHERE room_id = $1 ORDER BY timestamp ASC`
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

	return messages, rows.Err()
}

// GetMessagesSince retrieves messages for a specific room since a given time.
func (s *DBStore) GetMessagesSince(roomID string, since time.Time) ([]*models.Message, error) {
	query := `SELECT id, room_id, sender_id, content, timestamp FROM messages WHERE room_id = ? AND timestamp > ? ORDER BY timestamp ASC`
	if !s.config.IsSQLite() {
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

	return messages, rows.Err()
}
