package ws

import (
	"main/token"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, you should implement proper origin checking
	},
}

// Handler handles WebSocket connections
type Handler struct {
	manager    *Manager
	tokenMaker token.Maker
}

// NewHandler creates a new WebSocket handler
func NewHandler(manager *Manager, tokenMaker token.Maker) *Handler {
	return &Handler{
		manager:    manager,
		tokenMaker: tokenMaker,
	}
}

// HandleConnection upgrades HTTP connection to WebSocket and handles the connection
func (h *Handler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// Authenticate the connection
	// authHeader := r.Header.Get("Authorization")

	// TODO: check if this is correct, mainly r.Context() gives required data

	log.Info().Msg("Authenticating WebSocket connection")

	// Log context information with specific auth details
	// log.Info().
	// 	Interface("headers", r.Header).
	// 	Interface("context", r.Context()).
	// 	Str("method", r.Method).
	// 	Str("url", r.URL.String()).
	// 	Str("remote_addr", r.RemoteAddr).
	// 	Msg("WebSocket connection context")

	// log.Info().Msg("Getting auth token")
	// log.Info().Msg(r.Header.Get("Auth_token"))

	payload, err := h.tokenMaker.AuthenticateUser(r.Header.Get("Auth_token"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to authenticate WebSocket connection")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade connection")
		return
	}

	// Create new client using the username from the token
	client := &Client{
		ID:      payload.Username, // Use username from token as client ID
		Conn:    conn,
		Manager: h.manager,
	}

	// Register client
	h.manager.register <- client

	// Start listening for messages
	go h.handleMessages(client)
}

// handleMessages processes incoming messages from the client
func (h *Handler) handleMessages(client *Client) {
	defer func() {
		h.manager.unregister <- client
	}()

	for {
		var message Message
		err := client.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Str("clientID", client.ID).Msg("Error reading message")
			}
			break
		}

		// Validate that the message playerID matches the authenticated user
		if message.PlayerID != client.ID {
			log.Error().
				Str("clientID", client.ID).
				Str("messagePlayerID", message.PlayerID).
				Msg("PlayerID in message does not match authenticated user")
			continue
		}

		// Process message based on type
		switch message.Type {
		case "create_game":
			gameID := message.GameID
			h.manager.CreateGame(gameID, client.ID)
			client.GameID = gameID
			h.broadcastGameState(gameID)

		case "join_game":
			gameID := message.GameID
			if h.manager.JoinGame(gameID, client.ID) {
				client.GameID = gameID
				h.broadcastGameState(gameID)
			}

		case "make_move":
			if data, ok := message.Data.(map[string]interface{}); ok {
				if position, ok := data["position"].(float64); ok {
					if h.manager.MakeMove(message.GameID, client.ID, int(position)) {
						h.broadcastGameState(message.GameID)
					}
				}
			}
		}
	}
}

// broadcastGameState sends the current game state to all players in the game
func (h *Handler) broadcastGameState(gameID string) {
	h.manager.mutex.RLock()
	game, exists := h.manager.games[gameID]
	h.manager.mutex.RUnlock()

	if !exists {
		return
	}

	h.manager.broadcast <- &Message{
		Type:   "game_state",
		GameID: gameID,
		Data:   game,
	}
}
