package internal

import (
	"net/http"
)

// NewRouter creates a new HTTP router with WebSocket endpoints
func NewRouter(hub *Hub) *http.ServeMux {
	router := http.NewServeMux()
	
	// Add a simple handler for the root path
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("WebSocket Chat Server"))
	})
	
	// Add WebSocket handler
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})
	
	return router
}
