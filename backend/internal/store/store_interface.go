package store

import (
	"backend/internal/models"
	"time"
)

// StoreInterface defines the methods that any store implementation must provide.

type StoreInterface interface {
	Migrate() error
	Close() error

	// User methods
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)

	// Room methods
	CreateRoom(room *models.ChatRoom) error
	GetRoomsByUserID(userID string) ([]*models.ChatRoom, error)
	AddRoomMember(member *models.RoomMember) error

	// Message methods
	SaveMessage(message *models.Message) error
	GetMessagesByRoom(roomID string) ([]*models.Message, error)
	GetMessagesSince(roomID string, since time.Time) ([]*models.Message, error)
}
