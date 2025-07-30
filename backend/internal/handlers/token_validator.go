package handlers

import (
	"backend/internal/models"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWT secret key - must match the one in middleware/auth.go
const jwtSecret = "your-256-bit-secret-key-change-this-in-production"

// validateToken validates a JWT token and returns the user if valid
func validateToken(tokenString string) (*models.User, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Create a user from the claims
	userID, _ := claims["sub"].(string)
	username, _ := claims["username"].(string)

	if userID == "" || username == "" {
		return nil, errors.New("invalid token: missing user information")
	}

	user := &models.User{
		ID:       userID,
		Username: username,
	}

	return user, nil
}
