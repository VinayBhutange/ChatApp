package services

import (
	"backend/internal/models"
	"backend/internal/store"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// JWT secret key used to sign tokens
const jwtSecret = "your-256-bit-secret-key-change-this-in-production"

// TokenExpiration defines how long a JWT token is valid
const TokenExpiration = 24 * time.Hour

// UserService provides user-related business logic.
type UserService struct {
	store store.StoreInterface
}

// TokenClaims represents the claims in the JWT token
type TokenClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewUserService creates a new UserService.
func NewUserService(s store.StoreInterface) *UserService {
	return &UserService{store: s}
}

// RegisterUser handles the business logic of creating a new user.
func (s *UserService) RegisterUser(username, password string) (*models.User, error) {
	log.Printf("RegisterUser: Starting registration for username: %s", username)
	
	// Skip the username check and directly try to create the user
	// The database will enforce uniqueness constraint on the username
	
	// Hash the password for security
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("RegisterUser: Failed to hash password: %v", err)
		return nil, err
	}

	// Create a new user model
	userID := uuid.NewString()
	log.Printf("RegisterUser: Generated UUID: %s", userID)
	newUser := &models.User{
		ID:       userID,
		Username: username,
		Password: string(hashedPassword),
	}

	// Save the user to the database
	log.Printf("RegisterUser: Attempting to save user to database")
	if err = s.store.CreateUser(newUser); err != nil {
		log.Printf("RegisterUser: Database error creating user: %v", err)
		
		// Check if the error is due to a duplicate username
		if err.Error() == "UNIQUE constraint failed: users.username" || 
		   err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			return nil, errors.New("username already exists")
		}
		
		return nil, err
	}

	log.Printf("RegisterUser: User successfully created in database")
	return newUser, nil
}

// GetUserByUsername retrieves a user by username.
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.store.GetUserByUsername(username)
}

// AuthenticateUser validates a username and password, returning a JWT token if valid.
func (s *UserService) AuthenticateUser(username, password string) (string, error) {
	// Get the user from the database
	user, err := s.store.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Additional safety check to prevent nil pointer dereference
	if user == nil {
		return "", errors.New("invalid username or password")
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Generate a JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateJWT creates a new JWT token for a user
func (s *UserService) generateJWT(user *models.User) (string, error) {
	// Create the claims
	claims := TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "chat-app",
			Subject:   user.ID,
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
