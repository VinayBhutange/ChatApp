package handlers

import (
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"net/http"
)

// RoomHandler handles HTTP requests for chat room-related actions.
type RoomHandler struct {
	roomService *services.RoomService
}

// NewRoomHandler creates a new RoomHandler.
func NewRoomHandler(roomService *services.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

// CreateRoomRequest defines the expected JSON body for a room creation request.
type CreateRoomRequest struct {
	Name string `json:"name"`
}

// GetRoomsResponse defines the JSON response for listing rooms.
type GetRoomsResponse struct {
	Rooms []models.ChatRoom `json:"rooms"`
}

// CreateRoom handles the creation of a new chat room.
func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Room name is required", http.StatusBadRequest)
		return
	}

	// Create all new rooms as 'private' by default
	room, err := h.roomService.CreateRoom(req.Name, user.ID, "private")
	if err != nil {
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(room)
}

// GetRooms handles listing all rooms a user has access to.
func (h *RoomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rooms, err := h.roomService.GetRoomsForUser(user.ID)
	if err != nil {
		http.Error(w, "Failed to retrieve rooms", http.StatusInternalServerError)
		return
	}

	response := GetRoomsResponse{
		Rooms: rooms,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
