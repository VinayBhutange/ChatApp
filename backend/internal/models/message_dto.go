package models

import "time"

// MessageDTO is the data transfer object for a message sent over WebSocket.
// It includes the sender's username for easy display on the frontend.
type MessageDTO struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"roomId"`
	SenderID  string    `json:"senderId"`
	Sender    string    `json:"sender"` // Username of the sender
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
