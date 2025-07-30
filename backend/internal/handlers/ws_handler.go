package handlers

import (
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections for real-time chat.
type WebSocketHandler struct {
	messageService *services.MessageService
	clients        map[*Client]bool
	broadcast      chan *models.Message
	register       chan *Client
	unregister     chan *Client
}

// Client represents a connected WebSocket client.
type Client struct {
	hub      *WebSocketHandler
	conn     *websocket.Conn
	send     chan []byte
	user     *models.User
	roomID   string
}

// NewWebSocketHandler creates a new WebSocketHandler.
func NewWebSocketHandler(messageService *services.MessageService) *WebSocketHandler {
	return &WebSocketHandler{
		messageService: messageService,
		clients:        make(map[*Client]bool),
		broadcast:      make(chan *models.Message),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
	}
}

// Run starts the WebSocket hub.
func (h *WebSocketHandler) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client connected: %s in room %s", client.user.Username, client.roomID)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected: %s", client.user.Username)
			}

		case message := <-h.broadcast:
			// Save message to database
			err := h.messageService.SaveMessage(message)
			if err != nil {
				log.Printf("Error saving message: %v", err)
				continue
			}

			// Broadcast to all clients in the same room
			messageJSON, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			for client := range h.clients {
				if client.roomID == message.RoomID {
					select {
					case client.send <- messageJSON:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}

// ServeWs handles WebSocket requests from clients.
func (h *WebSocketHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	// Get room ID from query parameter
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}

	// Get authenticated user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		// For WebSocket connections, the token might be in the query string
		token := r.URL.Query().Get("token")
		if token != "" {
			// Validate the token manually
			user, err := validateToken(token)
			if err != nil {
				log.Printf("Invalid token in WebSocket connection: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// Continue with the authenticated user
			log.Printf("WebSocket authenticated with token from query: %s", user.Username)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// Allow all origins for development
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create new client
	client := &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		user:   user,
		roomID: roomID,
	}
	client.hub.register <- client

	// Start goroutines for reading and writing messages
	go client.writePump()
	go client.readPump()
}

// readPump pumps messages from the WebSocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(4096) // Max message size
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, msgBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		// Parse the message
		var msgData struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal(msgBytes, &msgData); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Create a new message
		message := &models.Message{
			ID:        services.GenerateUUID(),
			RoomID:    c.roomID,
			SenderID:  c.user.ID,
			Content:   msgData.Content,
			Timestamp: time.Now(),
		}

		// Send to broadcast channel
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the WebSocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second) // Send pings to client
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
