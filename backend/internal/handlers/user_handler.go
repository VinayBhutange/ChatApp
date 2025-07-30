package handlers

import (
	"backend/internal/services"
	"encoding/json"
	"log"
	"net/http"
)

// UserHandler handles HTTP requests for user-related actions.
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterRequest defines the expected JSON body for a registration request.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register handles the user registration endpoint.
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to register user: %s", req.Username)
	user, err := h.userService.RegisterUser(req.Username, req.Password)
	if err != nil {
		// Log the detailed error
		log.Printf("Failed to register user: %v", err)
		
		// Check for specific error types
		if err.Error() == "username already exists" {
			http.Error(w, "Username already exists. Please choose a different username.", http.StatusConflict) // 409 Conflict
		} else {
			http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	log.Printf("User registered successfully: %s (ID: %s)", user.Username, user.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// LoginRequest defines the expected JSON body for a login request.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse defines the JSON response for a successful login.
type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"user"`
}

// Login handles the user login endpoint.
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate the user
	token, err := h.userService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Get the user details from the service
	user, err := h.userService.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Failed to retrieve user details", http.StatusInternalServerError)
		return
	}

	// Create the response
	response := LoginResponse{
		Token: token,
	}
	response.User.ID = user.ID
	response.User.Username = user.Username

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
