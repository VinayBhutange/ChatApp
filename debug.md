# Chat Application Debug Guide

## Table of Contents
1. [Architecture Overview](#architecture-overview)
2. [Module Breakdown](#module-breakdown)
3. [API Endpoints](#api-endpoints)
4. [Database Operations](#database-operations)
5. [WebSocket Implementation](#websocket-implementation)
6. [Debugging Setup](#debugging-setup)
7. [Postman Testing Guide](#postman-testing-guide)
8. [Common Issues & Troubleshooting](#common-issues--troubleshooting)

## Architecture Overview

The chat application follows a clean architecture pattern with clear separation of concerns:

```
backend/
├── cmd/server/main.go          # Application entry point
├── internal/
│   ├── api/router.go           # REST API routes
│   ├── handlers/               # HTTP request handlers
│   ├── services/               # Business logic layer
│   ├── store/                  # Data access layer
│   ├── models/                 # Data models
│   ├── middleware/             # Authentication & CORS
│   ├── config/                 # Configuration management
│   ├── hub.go                  # WebSocket hub
│   ├── client.go               # WebSocket client
│   └── router.go               # WebSocket router
└── chatapp.db                  # SQLite database
```

## Module Breakdown

### 1. User Authentication Module

**Files:**
- `handlers/user_handler.go` - HTTP handlers for user operations
- `services/user_service.go` - User business logic
- `models/models.go` - User model definition

**Key Components:**
- User registration with password hashing
- JWT token-based authentication
- Login validation

**Debug Points:**
```go
// In user_service.go - Add logging
log.Printf("Registering user: %s", req.Username)
log.Printf("Password hash generated successfully")
log.Printf("User created with ID: %s", userID)
```

### 2. Chat Room Module

**Files:**
- `handlers/room_handler.go` - Room HTTP handlers
- `services/room_service.go` - Room business logic
- `models/models.go` - ChatRoom model

**Key Components:**
- Room creation with owner assignment
- Room listing for authenticated users
- Room membership management

**Debug Points:**
```go
// In room_service.go - Add logging
log.Printf("Creating room: %s for owner: %s", req.Name, ownerID)
log.Printf("Room created with ID: %s", roomID)
```

### 3. Real-time Messaging Module

**Files:**
- `handlers/ws_handler.go` - WebSocket connection handler
- `services/message_service.go` - Message business logic
- `hub.go` - WebSocket hub for managing connections
- `client.go` - Individual client connection

**Key Components:**
- WebSocket connection management
- Real-time message broadcasting
- Message persistence to database

**Debug Points:**
```go
// In ws_handler.go - Add logging
log.Printf("WebSocket connection established for user: %s", userID)
log.Printf("Message received: %+v", message)
log.Printf("Broadcasting message to room: %s", message.RoomID)
```

### 4. Database Module

**Files:**
- `store/db_store.go` - Database connection and migration
- `store/db_store_methods.go` - CRUD operations
- `store/store_interface.go` - Interface definition

**Key Components:**
- SQLite database with migration support
- CRUD operations for all entities
- Connection pooling and error handling

## API Endpoints

### Authentication Endpoints

#### 1. User Registration
- **Endpoint:** `POST /api/register`
- **Purpose:** Create new user account
- **Request Body:**
```json
{
  "username": "testuser",
  "password": "password123"
}
```
- **Response:**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": "user-uuid",
    "username": "testuser"
  }
}
```

#### 2. User Login
- **Endpoint:** `POST /api/login`
- **Purpose:** Authenticate user and get JWT token
- **Request Body:**
```json
{
  "username": "testuser",
  "password": "password123"
}
```
- **Response:**
```json
{
  "token": "jwt-token-here",
  "user": {
    "id": "user-uuid",
    "username": "testuser"
  }
}
```

### Room Management Endpoints

#### 3. Create Room
- **Endpoint:** `POST /api/rooms/create`
- **Purpose:** Create new chat room
- **Headers:** `Authorization: Bearer <jwt-token>`
- **Request Body:**
```json
{
  "name": "General Chat",
  "roomType": "public"
}
```
- **Response:**
```json
{
  "message": "Room created successfully",
  "room": {
    "id": "room-uuid",
    "name": "General Chat",
    "ownerId": "user-uuid",
    "roomType": "public"
  }
}
```

#### 4. List Rooms
- **Endpoint:** `GET /api/rooms`
- **Purpose:** Get all available rooms
- **Headers:** `Authorization: Bearer <jwt-token>`
- **Response:**
```json
{
  "rooms": [
    {
      "id": "room-uuid",
      "name": "General Chat",
      "ownerId": "user-uuid",
      "roomType": "public"
    }
  ]
}
```

### WebSocket Endpoint

#### 5. WebSocket Connection
- **Endpoint:** `WS /api/ws`
- **Purpose:** Establish real-time connection
- **Query Parameters:** `token=<jwt-token>&roomId=<room-uuid>`

## Database Operations

### Database Schema

```sql
-- Users table
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

-- Chat rooms table
CREATE TABLE chat_rooms (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    owner_id TEXT NOT NULL,
    room_type TEXT NOT NULL DEFAULT 'public',
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

-- Messages table
CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    room_id TEXT NOT NULL,
    sender_id TEXT NOT NULL,
    content TEXT NOT NULL,
    timestamp DATETIME NOT NULL,
    FOREIGN KEY (room_id) REFERENCES chat_rooms(id),
    FOREIGN KEY (sender_id) REFERENCES users(id)
);

-- Room members table
CREATE TABLE room_members (
    room_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    PRIMARY KEY (room_id, user_id),
    FOREIGN KEY (room_id) REFERENCES chat_rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### Key Database Operations

1. **User Operations:**
   - `CreateUser()` - Insert new user with hashed password
   - `GetUserByUsername()` - Retrieve user for authentication
   - `GetUserByID()` - Get user details by ID

2. **Room Operations:**
   - `CreateRoom()` - Create new chat room
   - `GetRooms()` - List all rooms
   - `GetRoomByID()` - Get specific room details

3. **Message Operations:**
   - `CreateMessage()` - Store new message
   - `GetMessagesByRoom()` - Retrieve room message history

## WebSocket Implementation

### Connection Flow
1. Client connects to `/api/ws` with JWT token and room ID
2. Server validates token and extracts user information
3. Client is registered with the WebSocket hub
4. Client joins specified room for message broadcasting

### Message Flow
1. Client sends message through WebSocket connection
2. Server validates message and user permissions
3. Message is stored in database
4. Message is broadcast to all clients in the room

### Hub Architecture
- **Hub:** Central manager for all WebSocket connections
- **Client:** Individual connection wrapper with send/receive channels
- **Broadcast:** Message distribution to room members

## Debugging Setup

### 1. Enable Detailed Logging

Add comprehensive logging throughout the application:

```go
// In main.go
log.SetFlags(log.LstdFlags | log.Lshortfile)

// In handlers
log.Printf("[DEBUG] Request received: %s %s", r.Method, r.URL.Path)
log.Printf("[DEBUG] Request body: %s", string(body))
log.Printf("[DEBUG] Response: %+v", response)

// In services
log.Printf("[DEBUG] Service method called with params: %+v", params)
log.Printf("[DEBUG] Database query result: %+v", result)

// In WebSocket
log.Printf("[DEBUG] WebSocket message: %+v", message)
log.Printf("[DEBUG] Active connections: %d", len(hub.clients))
```

### 2. Database Query Logging

Enable SQL query logging:

```go
// In db_store.go
func (s *DBStore) logQuery(query string, args ...interface{}) {
    log.Printf("[SQL] Query: %s", query)
    log.Printf("[SQL] Args: %+v", args)
}
```

### 3. Error Handling Enhancement

Improve error messages with context:

```go
func (s *DBStore) CreateUser(user *models.User) error {
    query := `INSERT INTO users (id, username, password) VALUES (?, ?, ?)`
    _, err := s.db.Exec(query, user.ID, user.Username, user.Password)
    if err != nil {
        log.Printf("[ERROR] Failed to create user %s: %v", user.Username, err)
        return fmt.Errorf("database error creating user %s: %w", user.Username, err)
    }
    log.Printf("[DEBUG] User created successfully: %s", user.Username)
    return nil
}
```

## Postman Testing Guide

### Setup Postman Environment

1. Create a new environment in Postman
2. Add variables:
   - `base_url`: `http://localhost:8082`
   - `auth_token`: (will be set after login)
   - `user_id`: (will be set after login)
   - `room_id`: (will be set after room creation)

### Test Sequence

#### Step 1: Test Server Health
```
GET {{base_url}}/
Expected: "WebSocket Chat Server" or 200 OK
```

#### Step 2: Register New User
```
POST {{base_url}}/api/register
Content-Type: application/json

{
  "username": "testuser1",
  "password": "password123"
}

Expected Response:
{
  "message": "User registered successfully",
  "user": {
    "id": "generated-uuid",
    "username": "testuser1"
  }
}
```

**Postman Test Script:**
```javascript
pm.test("User registration successful", function () {
    pm.response.to.have.status(200);
    const response = pm.response.json();
    pm.expect(response.message).to.include("successfully");
    pm.environment.set("user_id", response.user.id);
});
```

#### Step 3: Login User
```
POST {{base_url}}/api/login
Content-Type: application/json

{
  "username": "testuser1",
  "password": "password123"
}

Expected Response:
{
  "token": "jwt-token-string",
  "user": {
    "id": "user-uuid",
    "username": "testuser1"
  }
}
```

**Postman Test Script:**
```javascript
pm.test("Login successful", function () {
    pm.response.to.have.status(200);
    const response = pm.response.json();
    pm.expect(response.token).to.exist;
    pm.environment.set("auth_token", response.token);
});
```

#### Step 4: Create Chat Room
```
POST {{base_url}}/api/rooms/create
Content-Type: application/json
Authorization: Bearer {{auth_token}}

{
  "name": "Test Room",
  "roomType": "public"
}

Expected Response:
{
  "message": "Room created successfully",
  "room": {
    "id": "room-uuid",
    "name": "Test Room",
    "ownerId": "user-uuid",
    "roomType": "public"
  }
}
```

**Postman Test Script:**
```javascript
pm.test("Room creation successful", function () {
    pm.response.to.have.status(200);
    const response = pm.response.json();
    pm.expect(response.room.id).to.exist;
    pm.environment.set("room_id", response.room.id);
});
```

#### Step 5: List Rooms
```
GET {{base_url}}/api/rooms
Authorization: Bearer {{auth_token}}

Expected Response:
{
  "rooms": [
    {
      "id": "room-uuid",
      "name": "Test Room",
      "ownerId": "user-uuid",
      "roomType": "public"
    }
  ]
}
```

#### Step 6: Test WebSocket Connection

For WebSocket testing, use a WebSocket client or browser console:

```javascript
// In browser console
const token = "your-jwt-token";
const roomId = "your-room-id";
const ws = new WebSocket(`ws://localhost:8082/api/ws?token=${token}&roomId=${roomId}`);

ws.onopen = function() {
    console.log("WebSocket connected");
    
    // Send a test message
    ws.send(JSON.stringify({
        type: "message",
        content: "Hello, World!",
        roomId: roomId
    }));
};

ws.onmessage = function(event) {
    console.log("Message received:", JSON.parse(event.data));
};

ws.onerror = function(error) {
    console.error("WebSocket error:", error);
};
```

### Advanced Postman Tests

#### Test Authentication Failure
```
POST {{base_url}}/api/rooms/create
Content-Type: application/json
Authorization: Bearer invalid-token

Expected: 401 Unauthorized
```

#### Test Duplicate Username Registration
```
POST {{base_url}}/api/register
Content-Type: application/json

{
  "username": "testuser1",
  "password": "different-password"
}

Expected: 400 Bad Request (username already exists)
```

#### Test Invalid Login
```
POST {{base_url}}/api/login
Content-Type: application/json

{
  "username": "testuser1",
  "password": "wrong-password"
}

Expected: 401 Unauthorized
```

## Common Issues & Troubleshooting

### 1. Database Issues

**Problem:** `database is locked` error
**Solution:**
```bash
# Stop the server and delete database files
rm chatapp.db chat.db
# Restart server to recreate database with proper schema
```

**Problem:** `no such column: owner_id`
**Solution:** Database schema mismatch - delete and recreate database

### 2. Authentication Issues

**Problem:** JWT token validation fails
**Debug Steps:**
1. Check token format in Authorization header
2. Verify token hasn't expired
3. Check JWT secret key consistency
4. Add logging in middleware:

```go
log.Printf("[DEBUG] Authorization header: %s", authHeader)
log.Printf("[DEBUG] Token validation result: %v", valid)
```

### 3. WebSocket Issues

**Problem:** WebSocket connection fails
**Debug Steps:**
1. Check if JWT token is valid
2. Verify room ID exists
3. Check CORS settings
4. Add WebSocket connection logging:

```go
log.Printf("[DEBUG] WebSocket upgrade request from: %s", r.RemoteAddr)
log.Printf("[DEBUG] WebSocket query params: %v", r.URL.Query())
```

### 4. CORS Issues

**Problem:** Frontend can't connect to backend
**Solution:** Ensure CORS middleware is properly configured:

```go
// In middleware/cors.go
w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
```

### 5. Port Conflicts

**Problem:** Server fails to start on port 8082
**Debug Steps:**
1. Check if port is already in use: `netstat -an | findstr 8082`
2. Kill existing process or use different port
3. Set PORT environment variable: `set PORT=8083`

## Running the Application

### Start Backend Server
```bash
cd backend
go run cmd/server/main.go
```

### Expected Output
```
2024/08/08 20:14:53 Starting Chat Application Server...
2024/08/08 20:14:53 Using SQLite database
2024/08/08 20:14:53 Database connection established successfully
2024/08/08 20:14:53 Starting database migration
2024/08/08 20:14:53 Database migration completed successfully
2024/08/08 20:14:53 WebSocket hub started
2024/08/08 20:14:53 Server starting on port 8082...
```

### Verify Server is Running
```bash
curl http://localhost:8082/
# Expected: "WebSocket Chat Server"
```

This debug guide provides comprehensive information for understanding, testing, and troubleshooting the chat application backend. Use the Postman collection to systematically test each module and refer to the debugging sections when issues arise.
