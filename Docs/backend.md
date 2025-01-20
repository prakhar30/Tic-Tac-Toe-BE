# Game Lobby Backend Plan

## Database Schema

### Users Table
- `id`: UUID
- `username`: String
- `email`: String
- `password_hash`: String
- `created_at`: Timestamp

### Games Table
- `id`: UUID
- `name`: String
- `host_user_id`: UUID (Foreign Key to Users)
- `status`: Enum (e.g., "waiting", "in_progress", "completed")
- `current_state`: String (to store the board state, e.g., "XOX O O  ")
- `next_turn_user_id`: UUID (Foreign Key to Users, indicating whose turn it is)

### GameParticipants Table
- `id`: UUID
- `game_id`: UUID (Foreign Key to Games)
- `user_id`: UUID (Foreign Key to Users)
- `joined_at`: Timestamp

## API Endpoints (gRPC Services)

### User Service
- `RegisterUser`: Registers a new user.
- `LoginUser`: Authenticates a user and returns a token.

### Game Service
- `CreateGame`: Allows a user to create a new game.
- `JoinGame`: Allows a user to join an existing game.
- `StartGame`: Starts a game when all participants are ready.
- `MakeMove`: Allows a player to make a move, updates the game state, and switches turns.
- `GetGameStatus`: Retrieves the current status and state of a game.

### Lobby Service
- `ListAvailableGames`: Lists all games that are in "waiting" status.
- `GetGameParticipants`: Lists all participants in a specific game.

## WebSockets for Real-Time Communication

### Overview
WebSockets will be used to enable real-time communication between the client and server, allowing players to see updates to the game state immediately after a move is made.

### Implementation Plan

1. **Backend Setup:**
   - Use a WebSocket library in Go (e.g., `gorilla/websocket`) to handle WebSocket connections.
   - Create a WebSocket endpoint that clients can connect to.
   - Maintain a list of connected clients for each game.

2. **Handling Moves:**
   - When a player makes a move, update the game state in the database.
   - Broadcast the updated game state to all connected clients for that game.

3. **Frontend Setup:**
   - Establish a WebSocket connection to the server when the game page loads.
   - Listen for messages from the server and update the game state accordingly.


## Authentication
- Use JWT (JSON Web Tokens) for user authentication.
- Secure endpoints with middleware to check for valid tokens.

## gRPC Server Setup
- Define proto files for each service.
- Implement server logic in Go.
- Use gRPC for communication between client and server.

## Deployment Considerations
- Use Docker for containerization.
- Consider using Kubernetes for orchestration if scaling is needed.
- Use a cloud provider like AWS or GCP for hosting. 