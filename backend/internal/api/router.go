package api

import (
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/store"
	"net/http"
)

// NewRouter creates the main API router and registers all the application's routes.
func NewRouter(userHandler *handlers.UserHandler, roomHandler *handlers.RoomHandler, wsHandler *handlers.WebSocketHandler, dbStore store.StoreInterface) http.Handler {
	// Create test handler for debugging
	testHandler := handlers.NewTestHandler(dbStore)
	router := http.NewServeMux()

	// Public routes - no authentication required
	router.HandleFunc("/api/register", userHandler.Register)
	router.HandleFunc("/api/login", userHandler.Login)
	router.Handle("/api/rooms", middleware.RequireAuth(roomHandler.GetRooms)) // Protected endpoint to list rooms
	
	// Test endpoint for debugging registration issues
	router.HandleFunc("/api/test/register", testHandler.TestRegister)

	// Protected routes - require authentication
	router.Handle("/api/rooms/create", middleware.RequireAuth(roomHandler.CreateRoom))
		// The WebSocket handler performs its own authentication, so we don't need the RequireAuth middleware here.
	router.HandleFunc("/api/ws", wsHandler.ServeWs)

	// Apply CORS middleware to all routes
	return middleware.CORS(router)
}
