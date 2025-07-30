package services

import (
	"backend/internal/models"
	"backend/internal/store"

	"github.com/google/uuid"
)

// RoomService provides room-related business logic.
type RoomService struct {
	store store.StoreInterface
}

// NewRoomService creates a new RoomService.
func NewRoomService(s store.StoreInterface) *RoomService {
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

// GetRooms returns all chat rooms.
func (s *RoomService) GetRooms() ([]models.ChatRoom, error) {
	rooms, err := s.store.GetAllRooms()
	if err != nil {
		return nil, err
	}
	
	// Convert []*models.ChatRoom to []models.ChatRoom
	result := make([]models.ChatRoom, len(rooms))
	for i, room := range rooms {
		result[i] = *room
	}
	
	return result, nil
}
