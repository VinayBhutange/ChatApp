package handlers

import (
	"backend/internal/models"
	"backend/internal/store"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// DirectRegisterHandler provides a simplified registration endpoint that bypasses the service layer
type DirectRegisterHandler struct {
	store store.StoreInterface
}

// NewDirectRegisterHandler creates a new DirectRegisterHandler
func NewDirectRegisterHandler(store store.StoreInterface) *DirectRegisterHandler {
	return &DirectRegisterHandler{store: store}
}

// Register handles direct user registration
func (h *DirectRegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("DirectRegister: Received request")
	
	// Set CORS headers for this endpoint
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("DirectRegister: Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	log.Printf("DirectRegister: Attempting to register user: %s", req.Username)
	
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("DirectRegister: Failed to hash password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Create a new user
	userID := uuid.NewString()
	user := &models.User{
		ID:       userID,
		Username: req.Username,
		Password: string(hashedPassword),
	}
	
	// Save the user directly to the database
	if err := h.store.CreateUser(user); err != nil {
		log.Printf("DirectRegister: Failed to create user: %v", err)
		if err.Error() == "UNIQUE constraint failed: users.username" || 
		   err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			http.Error(w, "Username already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}
		return
	}
	
	// Return the created user (without password)
	response := struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}{
		ID:       user.ID,
		Username: user.Username,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	
	log.Printf("DirectRegister: Successfully registered user: %s with ID: %s", user.Username, user.ID)
}
