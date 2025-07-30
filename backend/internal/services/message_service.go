package services

import (
	"backend/internal/models"
	"backend/internal/store"
	"time"

	"github.com/google/uuid"
)

// MessageService provides message-related business logic.
type MessageService struct {
	store *store.Store
}

// NewMessageService creates a new MessageService.
func NewMessageService(s *store.Store) *MessageService {
	return &MessageService{store: s}
}

// GenerateUUID is a helper function to generate a UUID.
func GenerateUUID() string {
	return uuid.NewString()
}

// SaveMessage saves a message to the database.
func (s *MessageService) SaveMessage(message *models.Message) error {
	return s.store.SaveMessage(message)
}

// GetMessagesByRoom retrieves all messages for a specific room.
func (s *MessageService) GetMessagesByRoom(roomID string, limit, offset int) ([]models.Message, error) {
	return s.store.GetMessagesByRoom(roomID, limit, offset)
}

// GetRecentMessages retrieves the most recent messages for a specific room.
func (s *MessageService) GetRecentMessages(roomID string, limit int) ([]models.Message, error) {
	return s.store.GetMessagesByRoom(roomID, limit, 0)
}

// GetMessagesSince retrieves all messages for a room since a specific time.
func (s *MessageService) GetMessagesSince(roomID string, since time.Time) ([]models.Message, error) {
	return s.store.GetMessagesSince(roomID, since)
}
