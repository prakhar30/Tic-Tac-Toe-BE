# Tic-Tac-Toe Backend

A modern, secure, and scalable backend implementation for a real-time multiplayer Tic-Tac-Toe game, built with Go and following industry best practices.

## 🌟 Features

### Authentication & Security
- Token-based authentication using PASETO (Platform-Agnostic Security Tokens)
- Secure WebSocket connections with token validation
- Rate limiting and connection management
- SQL injection protection with prepared statements
- Environment-based configuration management

### Real-time Gameplay
- WebSocket-based real-time game state updates
- Support for multiple concurrent games
- Automatic game state synchronization
- Turn-based gameplay enforcement
- Win condition detection (horizontal, vertical, diagonal)

### API Design
- Clean architecture with separation of concerns
- gRPC API with protocol buffers
- WebSocket API for real-time communication
- Comprehensive error handling and validation
- Structured logging with zerolog

### Database & State Management
- PostgreSQL database with migrations
- Connection pooling for optimal performance
- Transactional data consistency
- In-memory game state management
- Concurrent access handling with mutexes

## 🛠 Technical Stack

- **IDE**: Cursor IDE
- **Language**: Go 1.22+
- **Database**: PostgreSQL
- **API Protocols**: 
  - gRPC with Protocol Buffers
  - WebSocket (gorilla/websocket)
- **Authentication**: PASETO tokens
- **Logging**: zerolog
- **Configuration**: Viper
- **Database Migration**: golang-migrate
- **Connection Pooling**: pgxpool
- **Error Handling**: Custom error types with proper propagation
- **Containerization**: Docker & Docker Compose

## 🏗 Architecture

```
├── api/            # API protocol definitions
├── db/             # Database migrations and queries
├── gapi/           # gRPC service implementations
├── pb/             # Generated Protocol Buffer code
├── token/          # Token management and authentication
├── utils/          # Utility functions and configurations
└── ws/             # WebSocket game logic and state management
```

### Key Components

1. **WebSocket Manager**: 
   - Handles real-time game state
   - Manages client connections
   - Broadcasts game updates
   - Implements game logic

2. **Authentication System**:
   - Token generation and validation
   - User session management
   - Secure WebSocket connections

3. **Game Logic**:
   - Turn management
   - Move validation
   - Win condition checking
   - Game state persistence

## 🚀 Getting Started

### Prerequisites
- Go 1.22 or higher
- PostgreSQL
- Docker & Docker Compose (optional)

### Running Locally

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/tic-tac-toe-be
   cd tic-tac-toe-be
   ```

2. Set up environment variables:
   ```bash
   cp app.env.example app.env
   # Edit app.env with your configuration
   ```

3. Start the database:
   ```bash
   docker-compose up -d postgres
   ```

4. Run database migrations:
   ```bash
   make migrateup
   ```

5. Start the server:
   ```bash
   make server
   ```

### Using Docker

```bash
docker-compose up --build
```

## 📡 API Documentation

### gRPC Endpoints
- `CreateUser`: Register new users
- `LoginUser`: Authenticate and receive tokens
- `UpdateUser`: Update user information
- `ValidateToken`: Verify token validity

### WebSocket Events
- `create_game`: Initialize a new game
- `join_game`: Join an existing game
- `make_move`: Make a move in the game
- `game_state`: Receive game state updates

## 🔒 Security Features

1. **Token-based Authentication**:
   - PASETO tokens for secure authentication
   - Token expiration and refresh mechanism
   - Secure token validation

2. **WebSocket Security**:
   - Authenticated connections
   - Rate limiting
   - Origin validation
   - Input sanitization

3. **Database Security**:
   - Prepared statements
   - Connection pooling
   - Transaction management
   - Secure credential handling

## 📈 Performance Considerations

- Connection pooling for database optimization
- Efficient game state management
- Concurrent request handling
- Memory-efficient data structures
- Proper error handling and resource cleanup

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
