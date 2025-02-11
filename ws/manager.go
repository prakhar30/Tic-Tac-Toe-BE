package ws

import (
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

// Message represents the WebSocket message structure
type Message struct {
	Type   string      `json:"type"`
	GameID string      `json:"gameId"`
	Data   interface{} `json:"data"`
}

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
			log.Info().Str("clientID", client.ID).Msg("New client connected")

		case client := <-m.unregister:
			m.mutex.Lock()
			if _, ok := m.clients[client.ID]; ok {
				delete(m.clients, client.ID)
				client.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				client.Conn.Close()
			}
			m.mutex.Unlock()
			log.Info().Str("clientID", client.ID).Msg("Client disconnected")

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
				client.Conn.Close()
				m.unregister <- client
			}
		}
	}
}

// CreateGame initializes a new game
func (m *Manager) CreateGame(gameID string, playerID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.games[gameID] = &GameState{
		Board:     [9]string{},
		Players:   make(map[string]string),
		Turn:      playerID,
		GameOver:  false,
		GameReady: false,
	}
	m.games[gameID].Players[playerID] = "X" // First player is X
}

// JoinGame adds a player to an existing game
func (m *Manager) JoinGame(gameID string, playerID string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if game, exists := m.games[gameID]; exists {
		if len(game.Players) < 2 {
			game.Players[playerID] = "O" // Second player is O
			game.GameReady = true
			return true
		}
	}
	return false
}

// MakeMove handles a player's move
func (m *Manager) MakeMove(gameID string, playerID string, position int) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	game, exists := m.games[gameID]
	if !game.GameReady || !exists || game.GameOver || game.Turn != playerID || position < 0 || position > 8 || game.Board[position] != "" {
		log.Error().Bool("gameReady", game.GameReady).Bool("exists", exists).Bool("gameOver", game.GameOver).Str("gameID", gameID).Str("playerID", playerID).Int("position", position).Msg("Invalid move")
		return false
	}

	game.Board[position] = game.Players[playerID]

	// Check for winner
	if winner := m.checkWinner(game.Board); winner != "" {
		game.Winner = playerID
		game.GameOver = true
	} else if m.isBoardFull(game.Board) {
		game.GameOver = true
	} else {
		// Switch turns
		for pid := range game.Players {
			if pid != playerID {
				game.Turn = pid
				break
			}
		}
	}

	return true
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
