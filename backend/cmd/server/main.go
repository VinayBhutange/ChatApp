package main

import (
	"backend/internal/api"
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/services"
	"backend/internal/store"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting Chat Application Server...")

	// Initialize database configuration
	dbConfig := config.NewDatabaseConfig()
	
	// Initialize store
	dbStore, err := store.NewDBStore(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbStore.Close()

	// Run database migrations
	if err := dbStore.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize services
	userService := services.NewUserService(dbStore)
	roomService := services.NewRoomService(dbStore)
	messageService := services.NewMessageService(dbStore)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	roomHandler := handlers.NewRoomHandler(roomService)
	wsHandler := handlers.NewWebSocketHandler(messageService)

	// Initialize router
	router := api.NewRouter(userHandler, roomHandler, wsHandler)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Start the server
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
