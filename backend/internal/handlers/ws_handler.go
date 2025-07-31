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

			// Create a DTO to include the sender's username
			messageDTO := models.MessageDTO{
				ID:        message.ID,
				RoomID:    message.RoomID,
				SenderID:  message.SenderID,
				Sender:    message.SenderUsername, // This is the key change
				Content:   message.Content,
				Timestamp: message.Timestamp,
			}

			// Broadcast the DTO to all clients in the same room
			messageJSON, err := json.Marshal(messageDTO)
			if err != nil {
				log.Printf("Error marshaling message DTO: %v", err)
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
	log.Printf("WebSocket connection request received from %s", r.RemoteAddr)
	
	// Log headers for debugging
	log.Printf("Request headers: %v", r.Header)
	
	// Get room ID from query parameter
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		log.Printf("WebSocket connection rejected: missing room_id parameter")
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}
	log.Printf("WebSocket connection for room: %s", roomID)

	// Get authenticated user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		// For WebSocket connections, the token might be in the query string
		token := r.URL.Query().Get("token")
		if token != "" {
			log.Printf("WebSocket using token from query parameter (length: %d)", len(token))
			// Validate the token manually
			var err error
			user, err = validateToken(token)
			if err != nil {
				log.Printf("Invalid token in WebSocket connection: %v", err)
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}
			// Continue with the authenticated user
			log.Printf("WebSocket authenticated successfully for user: %s (ID: %s)", user.Username, user.ID)
		} else {
			log.Printf("WebSocket connection rejected: no token provided")
			http.Error(w, "Unauthorized: no token provided", http.StatusUnauthorized)
			return
		}
	} else {
		log.Printf("WebSocket authenticated from context for user: %s", user.Username)
	}

	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// Allow all origins for development
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			log.Printf("WebSocket origin: %s", origin)
			return true // Allow all origins in development
		},
		EnableCompression: true,
	}

	// Add response headers to help with CORS
	headers := w.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Access-Control-Allow-Credentials", "true")
	headers.Add("Access-Control-Allow-Headers", "content-type, authorization")

	log.Printf("Attempting to upgrade connection to WebSocket")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		http.Error(w, "Could not upgrade connection: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("WebSocket connection successfully established")
	

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
		log.Printf("Client %s disconnecting from room %s", c.user.Username, c.roomID)
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// Set read parameters
	c.conn.SetReadLimit(4096) // Max message size
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		log.Printf("Received pong from client %s", c.user.Username)
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Set close handler
	c.conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("Client %s closing connection: code=%d, reason=%s", c.user.Username, code, text)
		message := websocket.FormatCloseMessage(code, "")
		c.conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	log.Printf("Started read pump for client %s in room %s", c.user.Username, c.roomID)

	for {
		_, msgBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message from client %s: %v", c.user.Username, err)
			} else {
				log.Printf("Client %s connection closed: %v", c.user.Username, err)
			}
			break
		}
		
		log.Printf("Received %d bytes from client %s", len(msgBytes), c.user.Username)

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

		// Add sender's username to the message before broadcasting
		message.SenderUsername = c.user.Username

		// Send to broadcast channel
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the WebSocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second) // Send pings to client every 30 seconds
	defer func() {
		log.Printf("Stopping write pump for client %s in room %s", c.user.Username, c.roomID)
		ticker.Stop()
		c.conn.Close()
	}()

	log.Printf("Started write pump for client %s in room %s", c.user.Username, c.roomID)

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The hub closed the channel
				log.Printf("Hub closed channel for client %s, sending close message", c.user.Username)
				err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Printf("Error sending close message to client %s: %v", c.user.Username, err)
				}
				return
			}

			log.Printf("Sending %d bytes to client %s", len(message), c.user.Username)
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Error getting writer for client %s: %v", c.user.Username, err)
				return
			}
			
			_, err = w.Write(message)
			if err != nil {
				log.Printf("Error writing message to client %s: %v", c.user.Username, err)
				return
			}

			// Add queued messages
			n := len(c.send)
			if n > 0 {
				log.Printf("Sending %d additional queued messages to client %s", n, c.user.Username)
			}
			
			for i := 0; i < n; i++ {
				_, err := w.Write([]byte{'\n'})
				if err != nil {
					log.Printf("Error writing newline to client %s: %v", c.user.Username, err)
					return
				}
				
				_, err = w.Write(<-c.send)
				if err != nil {
					log.Printf("Error writing queued message to client %s: %v", c.user.Username, err)
					return
				}
			}

			if err := w.Close(); err != nil {
				log.Printf("Error closing writer for client %s: %v", c.user.Username, err)
				return
			}

		case <-ticker.C:
			log.Printf("Sending ping to client %s", c.user.Username)
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("Error sending ping to client %s: %v", c.user.Username, err)
				return
			}
			log.Printf("Ping sent successfully to client %s", c.user.Username)
		}
	}
}
