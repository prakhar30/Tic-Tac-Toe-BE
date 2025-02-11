package ws

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// GameState represents the current state of a Tic-Tac-Toe game
type GameState struct {
	Board     [9]string         `json:"board"`
	Players   map[string]string `json:"players"` // map[playerID]symbol (X or O)
	Turn      string            `json:"turn"`    // playerID whose turn it is
	Winner    string            `json:"winner"`  // playerID of winner, empty if no winner
	GameOver  bool              `json:"gameOver"`
	GameReady bool              `json:"gameReady"`
}

// Client represents a connected player
type Client struct {
	ID      string
	Conn    *websocket.Conn
	GameID  string
	Symbol  string
	Manager *Manager
}

// GameError represents a game-related error
type GameError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *GameError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Message represents the WebSocket message structure
type Message struct {
	Type   string      `json:"type"`
	GameID string      `json:"gameId"`
	Data   interface{} `json:"data"`
	Error  *GameError  `json:"error,omitempty"`
}

// Game error codes
const (
	ErrGameNotFound     = "GAME_NOT_FOUND"
	ErrGameFull         = "GAME_FULL"
	ErrInvalidMove      = "INVALID_MOVE"
	ErrNotPlayersTurn   = "NOT_PLAYERS_TURN"
	ErrGameNotReady     = "GAME_NOT_READY"
	ErrPositionOccupied = "POSITION_OCCUPIED"
)

// Manager handles WebSocket connections and game states
type Manager struct {
	games      map[string]*GameState
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	mutex      sync.RWMutex
}

// NewManager creates a new WebSocket manager
func NewManager() *Manager {
	return &Manager{
		games:      make(map[string]*GameState),
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

// Start begins listening for WebSocket events
func (m *Manager) Start() {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			m.clients[client.ID] = client
			m.mutex.Unlock()
			log.Info().
				Str("client_id", client.ID).
				Int("total_clients", len(m.clients)).
				Msg("Client registered with manager")

		case client := <-m.unregister:
			m.mutex.Lock()
			if _, ok := m.clients[client.ID]; ok {
				log.Info().
					Str("client_id", client.ID).
					Str("game_id", client.GameID).
					Int("remaining_clients", len(m.clients)-1).
					Msg("Unregistering client")

				delete(m.clients, client.ID)
				client.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				client.Conn.Close()
			}
			m.mutex.Unlock()

		case message := <-m.broadcast:
			m.broadcastToGame(message)
		}
	}
}

// broadcastToGame sends a message to all clients in a specific game
func (m *Manager) broadcastToGame(message *Message) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, client := range m.clients {
		if client.GameID == message.GameID {
			err := client.Conn.WriteJSON(message)
			if err != nil {
				log.Error().Err(err).Str("clientID", client.ID).Msg("Error broadcasting message")
				// Send error message to client before potential disconnect
				errorMsg := &Message{
					Type:  "error",
					Error: &GameError{Code: "BROADCAST_ERROR", Message: "Failed to send message"},
				}
				client.Conn.WriteJSON(errorMsg)

				// Only disconnect if it's a fatal error
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					client.Conn.Close()
					m.unregister <- client
				}
			}
		}
	}
}

// CreateGame initializes a new game
func (m *Manager) CreateGame(gameID string, playerID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	log.Info().
		Str("game_id", gameID).
		Str("player_id", playerID).
		Int("total_games", len(m.games)+1).
		Msg("Creating new game")

	m.games[gameID] = &GameState{
		Board:     [9]string{},
		Players:   make(map[string]string),
		Turn:      playerID,
		GameOver:  false,
		GameReady: false,
	}
	m.games[gameID].Players[playerID] = "X" // First player is X

	log.Debug().
		Str("game_id", gameID).
		Str("player_id", playerID).
		Interface("game_state", m.games[gameID]).
		Msg("Game created successfully")
}

// JoinGame adds a player to an existing game
func (m *Manager) JoinGame(gameID string, playerID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	game, exists := m.games[gameID]
	if !exists {
		log.Warn().
			Str("game_id", gameID).
			Str("player_id", playerID).
			Msg("Attempted to join non-existent game")
		return &GameError{Code: ErrGameNotFound, Message: "Game not found"}
	}

	if len(game.Players) >= 2 {
		log.Warn().
			Str("game_id", gameID).
			Str("player_id", playerID).
			Interface("existing_players", game.Players).
			Msg("Attempted to join full game")
		return &GameError{Code: ErrGameFull, Message: "Game is already full"}
	}

	game.Players[playerID] = "O" // Second player is O
	game.GameReady = true

	log.Info().
		Str("game_id", gameID).
		Str("player_id", playerID).
		Interface("game_state", game).
		Msg("Player joined game successfully")

	return nil
}

// MakeMove handles a player's move
func (m *Manager) MakeMove(gameID string, playerID string, position int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	game, exists := m.games[gameID]
	if !exists {
		log.Warn().
			Str("game_id", gameID).
			Str("player_id", playerID).
			Int("position", position).
			Msg("Attempted move in non-existent game")
		return &GameError{Code: ErrGameNotFound, Message: "Game not found"}
	}

	if !game.GameReady {
		log.Warn().
			Str("game_id", gameID).
			Str("player_id", playerID).
			Bool("game_ready", game.GameReady).
			Msg("Attempted move in game not ready")
		return &GameError{Code: ErrGameNotReady, Message: "Game is not ready to start"}
	}

	if game.GameOver {
		log.Warn().
			Str("game_id", gameID).
			Str("player_id", playerID).
			Str("winner", game.Winner).
			Msg("Attempted move in finished game")
		return &GameError{Code: ErrGameNotReady, Message: "Game is already over"}
	}

	if game.Turn != playerID {
		log.Warn().
			Str("game_id", gameID).
			Str("player_id", playerID).
			Str("current_turn", game.Turn).
			Msg("Player attempted move out of turn")
		return &GameError{Code: ErrNotPlayersTurn, Message: "Not your turn"}
	}

	if position < 0 || position > 8 {
		log.Warn().
			Str("game_id", gameID).
			Str("player_id", playerID).
			Int("position", position).
			Msg("Invalid board position")
		return &GameError{Code: ErrInvalidMove, Message: "Invalid position"}
	}

	if game.Board[position] != "" {
		log.Warn().
			Str("game_id", gameID).
			Str("player_id", playerID).
			Int("position", position).
			Str("existing_symbol", game.Board[position]).
			Msg("Position already occupied")
		return &GameError{Code: ErrPositionOccupied, Message: "Position already occupied"}
	}

	game.Board[position] = game.Players[playerID]

	log.Debug().
		Str("game_id", gameID).
		Str("player_id", playerID).
		Int("position", position).
		Interface("board", game.Board).
		Msg("Move completed")

	// Check for winner
	if winner := m.checkWinner(game.Board); winner != "" {
		game.Winner = playerID
		game.GameOver = true
		log.Info().
			Str("game_id", gameID).
			Str("winner", playerID).
			Interface("final_board", game.Board).
			Msg("Game won")
	} else if m.isBoardFull(game.Board) {
		game.GameOver = true
		log.Info().
			Str("game_id", gameID).
			Interface("final_board", game.Board).
			Msg("Game ended in draw")
	} else {
		// Switch turns
		for pid := range game.Players {
			if pid != playerID {
				game.Turn = pid
				log.Debug().
					Str("game_id", gameID).
					Str("next_turn", pid).
					Msg("Turn switched to next player")
				break
			}
		}
	}

	return nil
}

// checkWinner determines if there's a winner
func (m *Manager) checkWinner(board [9]string) string {
	// Winning combinations
	lines := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}

	for _, line := range lines {
		if board[line[0]] != "" &&
			board[line[0]] == board[line[1]] &&
			board[line[1]] == board[line[2]] {
			return board[line[0]]
		}
	}
	return ""
}

// isBoardFull checks if the board is full (draw)
func (m *Manager) isBoardFull(board [9]string) bool {
	for _, cell := range board {
		if cell == "" {
			return false
		}
	}
	return true
}
