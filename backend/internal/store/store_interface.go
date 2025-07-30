package store

import (
	"backend/internal/models"
	"time"
)

// StoreInterface defines the methods that any store implementation must provide
type StoreInterface interface {
	// User operations
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	
	// Room operations
	CreateRoom(room *models.ChatRoom) error
	GetAllRooms() ([]*models.ChatRoom, error)
	
	// Message operations
	SaveMessage(message *models.Message) error
	GetMessagesByRoom(roomID string) ([]*models.Message, error)
	GetMessagesSince(roomID string, since time.Time) ([]*models.Message, error)
	
	// Database operations
	Migrate() error
	Close() error
}
