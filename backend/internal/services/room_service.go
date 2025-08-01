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
func (s *RoomService) CreateRoom(name, ownerID, roomType string) (*models.ChatRoom, error) {
	// Create a new room model
	newRoom := &models.ChatRoom{
		ID:       uuid.NewString(),
		Name:     name,
		OwnerID:  ownerID,
		RoomType: roomType,
	}

	// Save the room to the database
	if err := s.store.CreateRoom(newRoom); err != nil {
		return nil, err
	}

	// Add the owner as the first member of the room
	firstMember := &models.RoomMember{
		RoomID: newRoom.ID,
		UserID: ownerID,
		Status: "member", // The owner is automatically a member
	}

	if err := s.store.AddRoomMember(firstMember); err != nil {
		// In a real-world app, we might want to roll back the room creation here.
		return nil, err
	}

	return newRoom, nil
}

// GetRoomsForUser returns all public rooms plus private rooms the user is a member of.
func (s *RoomService) GetRoomsForUser(userID string) ([]models.ChatRoom, error) {
	rooms, err := s.store.GetRoomsByUserID(userID)
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
