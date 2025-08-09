# Real-Time Chat Application

A production-ready real-time chat application with user authentication, room management, and WebSocket messaging.

## Technology Stack

### Backend
- **Language**: Go (version 1.18)
- **Architecture**: Clean architecture with layered packages
- **Database**: SQLite (local development) / PostgreSQL (containerized)
- **Authentication**: JWT tokens
- **Real-time Communication**: WebSockets (gorilla/websocket)

### Frontend
- **Framework**: React with TypeScript
- **State Management**: React Context API
- **Styling**: Modern CSS with responsive design
- **API Communication**: Fetch API with proxy configuration

## Features

- User registration and authentication
- Secure password hashing with bcrypt
- JWT-based authentication
- Chat room creation and management
- Real-time messaging with WebSockets
- Message persistence in database
- Responsive UI for desktop and mobile

## Running the Application

### Local Development

1. **Start the Backend Server:**
   ```bash
   cd backend
   go run ./cmd/server
   ```

2. **Start the Frontend Development Server:**
   ```bash
   cd frontend
   npm start
   ```

3. **Access the Application:**
   Open your web browser and navigate to `http://localhost:3000`

### Docker Deployment

#### Using Docker Compose (Local Build)

1. **Build and Start the Containers:**
   ```bash
   docker-compose up --build
   ```

2. **Access the Application:**
   Open your web browser and navigate to `http://localhost:3000`

3. **Stop the Containers:**
   ```bash
   docker-compose down
   ```

#### Using Pre-built Docker Image from GHCR

The backend Docker image is automatically built and published to GitHub Container Registry on every push to the main branch.

1. **Pull and Run the Backend Image:**
   ```bash
   # Pull the latest image
   docker pull ghcr.io/your-username/chatapp/chatapp-backend:latest
   
   # Run the backend container
   docker run -d \
     --name chatapp-backend \
     -p 8082:8082 \
     -v $(pwd)/data:/root/data \
     ghcr.io/your-username/chatapp/chatapp-backend:latest
   ```

2. **Available Tags:**
   - `latest` - Latest build from main branch
   - `main` - Latest build from main branch
   - `v1.0.0` - Specific version tags
   - `main-<commit-sha>` - Specific commit builds

3. **Environment Variables:**
   ```bash
   docker run -d \
     --name chatapp-backend \
     -p 8082:8082 \
     -e DB_PATH=/root/data/chatapp.db \
     -v $(pwd)/data:/root/data \
     ghcr.io/your-username/chatapp/chatapp-backend:latest
   ```

## Project Structure

```
├── backend/                # Go backend
│   ├── cmd/                # Entry points
│   │   └── server/         # API server
│   └── internal/           # Internal packages
│       ├── api/            # API routes
│       ├── config/         # Configuration
│       ├── handlers/       # HTTP handlers
│       ├── middleware/     # HTTP middleware
│       ├── models/         # Data models
│       ├── services/       # Business logic
│       └── store/          # Data access
├── frontend/              # React frontend
│   ├── public/            # Static files
│   └── src/               # Source code
│       ├── components/    # React components
│       ├── contexts/      # React contexts
│       ├── services/      # API services
│       ├── styles/        # CSS styles
│       └── types/         # TypeScript types
└── docker-compose.yml    # Docker configuration
```
