# ChatApp Backend Documentation

This document provides a comprehensive overview of the backend architecture, API endpoints, and WebSocket implementation for the ChatApp project.

## 1. Project Structure (Clean Architecture)

The backend follows a clean, layered architecture, which separates concerns and makes the codebase maintainable and scalable.

-   **/cmd/server**: The main entry point of the application. `main.go` here is responsible for:
    -   Loading configuration.
    -   Initializing the database connection (`SQLite`).
    -   Creating instances of the `store`, `services`, and `handlers`.
    -   Setting up the main router.
    -   Starting the HTTP server on port `8082`.

-   **/internal/api**: Contains the main router (`router.go`) which defines all the application's API endpoints and connects them to their respective handlers.

-   **/internal/models**: Defines the core data structures (structs) used throughout the application, such as `User`, `Room`, and `Message`.

-   **/internal/store**: The data access layer. It is responsible for all communication with the database. It abstracts all SQL queries, so the rest of the application doesn't need to know about the database schema.

-   **/internal/services**: Contains the core business logic. For example, the `UserService` handles password hashing and user creation logic, while the `MessageService` would handle saving messages.

-   **/internal/handlers**: This layer handles the incoming HTTP requests. It parses request data (like JSON bodies), calls the appropriate services to perform business logic, and formats the HTTP responses.

-   **/internal/middleware**: Contains HTTP middleware.
    -   `CORS`: Handles Cross-Origin Resource Sharing to allow the frontend (on port 3000) to communicate with the backend.
    -   `RequireAuth`: Protects routes by validating JWT tokens from the `Authorization` header.

-   **/internal/hub.go & client.go**: These files form the core of the real-time WebSocket system.
    -   `hub.go`: Manages the collection of all active WebSocket clients and broadcasts messages to them.
    -   `client.go`: Represents a single WebSocket connection and manages reading and writing messages for that specific client.

## 2. API Endpoints

The application exposes a set of RESTful API endpoints for user management and chat room operations.

-   `POST /api/register`: Creates a new user account.
    -   **Request Body**: `{ "username": "...", "password": "..." }`
    -   **Response**: Success message or error.

-   `POST /api/login`: Authenticates a user and returns a JWT token.
    -   **Request Body**: `{ "username": "...", "password": "..." }`
    -   **Response**: `{ "token": "..." }`

-   `GET /api/rooms`: (Public) Returns a list of all available chat rooms.
    -   **Response**: `[{ "id": "...", "name": "..." }, ...]`

-   `POST /api/rooms/create`: (Protected) Creates a new chat room.
    -   **Requires**: Valid JWT in `Authorization` header.
    -   **Request Body**: `{ "name": "..." }`
    -   **Response**: The newly created room object.

-   `GET /api/ws`: (WebSocket Upgrade) The endpoint for initiating a WebSocket connection.
    -   **Query Parameters**: `room_id` and `token`.
    -   This is not a standard REST endpoint but the entry point for real-time communication.

## 3. WebSocket Workflow (Real-Time Messaging)

The real-time functionality is the most complex part of the backend. Here’s a step-by-step breakdown of how it works:

1.  **The Hub Starts**: When the application starts, a single instance of the `Hub` is created and runs in its own goroutine. The Hub is the central controller for all WebSocket communication. It has channels to handle client registrations, un-registrations, and incoming messages to be broadcast.

2.  **Client Connection**:
    -   The frontend, after a user logs in and joins a room, attempts to connect to the `ws://.../api/ws?room_id=...&token=...` endpoint.
    -   The `ws_handler.go` receives this request. It does **not** use the `RequireAuth` middleware because the token is in the URL, not the header.
    -   The handler manually validates the JWT token from the query parameter.

3.  **Upgrading to WebSocket**:
    -   If the token is valid, the handler uses the `gorilla/websocket` library's `Upgrader` to upgrade the standard HTTP connection to a persistent WebSocket connection.
    -   The `CheckOrigin` function in the upgrader is configured to allow connections from the frontend's origin (`http://localhost:3000`).

4.  **Client Creation**:
    -   A new `Client` object is created for this connection. This object holds a reference to the WebSocket connection, the user's details, and the room they joined.
    -   This new `Client` is registered with the `Hub` by sending it to the `hub.register` channel.

5.  **Pumping Messages (Goroutines)**:
    -   For each client, two dedicated goroutines are started:
        -   `readPump`: This loop continuously listens for new messages coming from the client's WebSocket connection. When a message is received, it is sent to the `hub.broadcast` channel.
        -   `writePump`: This loop continuously listens for messages on the client's personal `send` channel. When a message arrives, it is written out to the client's WebSocket connection.
    -   This two-pump system prevents a slow client from blocking the entire application.

6.  **Broadcasting a Message**:
    -   When the `Hub` receives a message on its `broadcast` channel (from a client's `readPump`), it first saves the message to the database via the `MessageService`.
    -   It then creates a `MessageDTO` (Data Transfer Object) that includes the sender's username.
    -   Finally, the `Hub` iterates over all registered clients. For each client that is in the same room as the message's sender, the `Hub` sends the `MessageDTO` to that client's personal `send` channel.

7.  **Client Disconnection**:
    -   If a client closes their browser or the connection is lost, the `readPump` will error out.
    -   The `defer` block in `readPump` ensures the client is unregistered from the `Hub` (via the `hub.unregister` channel) and the WebSocket connection is closed cleanly.

This architecture ensures that messages are efficiently and safely broadcast to all relevant clients in real-time.