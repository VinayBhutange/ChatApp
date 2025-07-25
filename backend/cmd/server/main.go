package main

import (
	"backend/internal/api"
	"backend/internal/handlers"
	"backend/internal/services"
	"backend/internal/store"
	"log"
	"net/http"
)

func main() {
	// Initialize the database store
	db, err := store.New("chat.db")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Run database migrations
	if err := db.Migrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database initialized and migrated successfully.")

	// --- Dependency Injection ---
	// Create services
	userService := services.NewUserService(db)
	roomService := services.NewRoomService(db)

	// Create handlers
	userHandler := handlers.NewUserHandler(userService)
	roomHandler := handlers.NewRoomHandler(roomService)

	// Create the main API router
	router := api.NewRouter(userHandler, roomHandler)

	log.Println("API server starting on :8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
