package services

import (
	"backend/internal/models"
	"backend/internal/store"

	"github.com/google/uuid"
)

// RoomService provides chat room-related business logic.
type RoomService struct {
	store *store.Store
}

// NewRoomService creates a new RoomService.
func NewRoomService(s *store.Store) *RoomService {
	return &RoomService{store: s}
}

// CreateRoom handles the business logic of creating a new chat room.
func (s *RoomService) CreateRoom(name, creatorID string) (*models.ChatRoom, error) {
	// Create a new room model
	newRoom := &models.ChatRoom{
		ID:   uuid.NewString(),
		Name: name,
	}

	// Save the room to the database
	if err := s.store.CreateRoom(newRoom); err != nil {
		return nil, err
	}

	return newRoom, nil
}

// GetAllRooms retrieves all chat rooms from the database.
func (s *RoomService) GetAllRooms() ([]models.ChatRoom, error) {
	return s.store.GetAllRooms()
}
