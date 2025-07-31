package middleware

import (
	"net/http"
)

// CORS middleware adds Cross-Origin Resource Sharing headers to responses
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request for debugging
		// Set CORS headers to allow the specific frontend origin
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		// Set other CORS headers
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Allow credentials
		w.Header().Set("Access-Control-Max-Age", "86400") // Cache preflight requests for 24 hours

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
