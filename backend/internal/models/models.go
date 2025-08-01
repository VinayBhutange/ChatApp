package models

import "time"

// User represents a user in the system.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Password is never returned in JSON responses
}

// ChatRoom represents a chat room in the system.
type ChatRoom struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	OwnerID  string `json:"ownerId" db:"owner_id"`
	RoomType string `json:"roomType" db:"room_type"` // 'public' or 'private'
}

// Message represents a chat message in the system.
type Message struct {
	ID             string    `json:"id" db:"id"`
	RoomID         string    `json:"roomId" db:"room_id"`
	SenderID       string    `json:"senderId" db:"sender_id"`
	Content        string    `json:"content" db:"content"`
	Timestamp      time.Time `json:"timestamp" db:"timestamp"`
	SenderUsername string    `json:"-"` // This field is for internal use and not stored in the DB
}

// RoomMember represents the relationship between a user and a room.
type RoomMember struct {
	RoomID string `json:"roomId" db:"room_id"`
	UserID string `json:"userId" db:"user_id"`
	Status string `json:"status" db:"status"` // e.g., "member", "pending"
}
