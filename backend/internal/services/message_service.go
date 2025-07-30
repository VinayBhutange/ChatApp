package services

import (
	"backend/internal/models"
	"backend/internal/store"
	"time"

	"github.com/google/uuid"
)

// MessageService provides message-related business logic.
type MessageService struct {
	store store.StoreInterface
}

// NewMessageService creates a new MessageService.
func NewMessageService(s store.StoreInterface) *MessageService {
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
func (s *MessageService) GetMessagesByRoom(roomID string) ([]models.Message, error) {
	messages, err := s.store.GetMessagesByRoom(roomID)
	if err != nil {
		return nil, err
	}
	
	// Convert []*models.Message to []models.Message
	result := make([]models.Message, len(messages))
	for i, msg := range messages {
		result[i] = *msg
	}
	
	return result, nil
}

// GetRecentMessages retrieves the most recent messages for a specific room.
func (s *MessageService) GetRecentMessages(roomID string, limit int) ([]models.Message, error) {
	// Since our store implementation now returns all messages (limited to 50),
	// we'll just use GetMessagesByRoom and limit the results in memory if needed
	messages, err := s.GetMessagesByRoom(roomID)
	if err != nil {
		return nil, err
	}
	
	// Limit the results if necessary
	if limit > 0 && limit < len(messages) {
		return messages[:limit], nil
	}
	return messages, nil
}

// GetMessagesSince retrieves all messages for a room since a specific time.
func (s *MessageService) GetMessagesSince(roomID string, since time.Time) ([]models.Message, error) {
	messages, err := s.store.GetMessagesSince(roomID, since)
	if err != nil {
		return nil, err
	}
	
	// Convert []*models.Message to []models.Message
	result := make([]models.Message, len(messages))
	for i, msg := range messages {
		result[i] = *msg
	}
	
	return result, nil
}
