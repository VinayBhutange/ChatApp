package handlers

import (
	"backend/internal/models"
	"backend/internal/store"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// TestHandler provides test endpoints for debugging
type TestHandler struct{
	store store.StoreInterface
}

// NewTestHandler creates a new TestHandler
func NewTestHandler(store store.StoreInterface) *TestHandler {
	return &TestHandler{store: store}
}

// getStore returns the store interface
func (h *TestHandler) getStore() (store.StoreInterface, error) {
	if h.store == nil {
		return nil, fmt.Errorf("store not initialized")
	}
	return h.store, nil
}

// TestRegister is a simplified registration endpoint for testing
func (h *TestHandler) TestRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("TestRegister: Received request")
	
	// Log request details
	log.Printf("TestRegister: Method: %s, Content-Type: %s", r.Method, r.Header.Get("Content-Type"))
	
	if r.Method == http.MethodOptions {
		// Handle preflight request
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}
	
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("TestRegister: Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	log.Printf("TestRegister: Received registration request for username: %s", req.Username)
	
	// Get a reference to the store
	dbStore, err := h.getStore()
	if err != nil {
		log.Printf("TestRegister: Failed to get store: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("TestRegister: Failed to hash password: %v", err)
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
	
	// Save the user to the database
	if err := dbStore.CreateUser(user); err != nil {
		log.Printf("TestRegister: Failed to create user: %v", err)
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
	
	log.Printf("TestRegister: Successfully registered user in database: %s with ID: %s", user.Username, user.ID)
}
