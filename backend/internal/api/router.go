package api

import (
	"backend/internal/handlers"
	"backend/internal/middleware"
	"net/http"
)

// NewRouter creates the main API router and registers all the application's routes.
func NewRouter(userHandler *handlers.UserHandler, roomHandler *handlers.RoomHandler) *http.ServeMux {
	router := http.NewServeMux()

	// Public routes - no authentication required
	router.HandleFunc("/api/register", userHandler.Register)
	router.HandleFunc("/api/login", userHandler.Login)
	router.HandleFunc("/api/rooms", roomHandler.GetRooms) // Public endpoint to list rooms

	// Protected routes - require authentication
	router.Handle("/api/rooms/create", middleware.RequireAuth(roomHandler.CreateRoom))

	return router
}
