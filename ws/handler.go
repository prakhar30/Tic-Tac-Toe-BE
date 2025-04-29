package ws

import (
	"fmt"
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
	log.Info().
		Str("remote_addr", r.RemoteAddr).
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Msg("New WebSocket connection attempt")

	formattedToken := fmt.Sprintf("Bearer %s", r.URL.Query().Get("token"))
	payload, err := h.tokenMaker.AuthenticateUser(formattedToken)
	if err != nil {
		log.Error().
			Err(err).
			Str("remote_addr", r.RemoteAddr).
			Msg("Failed to authenticate WebSocket connection")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Info().
		Str("username", payload.Username).
		Str("remote_addr", r.RemoteAddr).
		Msg("WebSocket connection authenticated")

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().
			Err(err).
			Str("username", payload.Username).
			Str("remote_addr", r.RemoteAddr).
			Msg("Failed to upgrade connection to WebSocket")
		return
	}

	// Create new client using the username from the token
	client := &Client{
		ID:      payload.Username,
		Conn:    conn,
		Manager: h.manager,
	}

	log.Info().
		Str("client_id", client.ID).
		Str("remote_addr", r.RemoteAddr).
		Msg("New client registered")

	// Register client
	h.manager.register <- client

	// Start listening for messages
	go h.handleMessages(client)
}

// handleMessages processes incoming messages from the client
func (h *Handler) handleMessages(client *Client) {
	defer func() {
		log.Info().
			Str("client_id", client.ID).
			Str("game_id", client.GameID).
			Msg("Client disconnecting, cleaning up")
		h.manager.unregister <- client
	}()

	for {
		var message Message
		err := client.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().
					Err(err).
					Str("client_id", client.ID).
					Str("game_id", client.GameID).
					Msg("Unexpected WebSocket closure")
			} else {
				log.Info().
					Str("client_id", client.ID).
					Str("game_id", client.GameID).
					Msg("WebSocket connection closed normally")
			}
			break
		}

		log.Debug().
			Str("client_id", client.ID).
			Str("game_id", client.GameID).
			Str("message_type", message.Type).
			Interface("message_data", message.Data).
			Msg("Received WebSocket message")

		var response *Message

		// Process message based on type
		switch message.Type {
		case "create_game":
			gameID := message.GameID
			log.Info().
				Str("client_id", client.ID).
				Str("game_id", gameID).
				Msg("Creating new game")

			h.manager.CreateGame(gameID, client.ID)
			client.GameID = gameID
			h.broadcastGameState(gameID)

		case "join_game":
			gameID := message.GameID
			log.Info().
				Str("client_id", client.ID).
				Str("game_id", gameID).
				Msg("Attempting to join game")

			if err := h.manager.JoinGame(gameID, client.ID); err != nil {
				if gameErr, ok := err.(*GameError); ok {
					log.Warn().
						Str("client_id", client.ID).
						Str("game_id", gameID).
						Str("error_code", gameErr.Code).
						Str("error_message", gameErr.Message).
						Msg("Failed to join game")

					response = &Message{
						Type:   "error",
						GameID: gameID,
						Error:  gameErr,
					}
					client.Conn.WriteJSON(response)
					continue
				}
			}

			log.Info().
				Str("client_id", client.ID).
				Str("game_id", gameID).
				Msg("Successfully joined game")

			client.GameID = gameID
			h.broadcastGameState(gameID)

		case "make_move":
			if data, ok := message.Data.(map[string]interface{}); ok {
				if position, ok := data["position"].(float64); ok {
					log.Info().
						Str("client_id", client.ID).
						Str("game_id", message.GameID).
						Float64("position", position).
						Msg("Attempting to make move")

					if err := h.manager.MakeMove(message.GameID, client.ID, int(position)); err != nil {
						if gameErr, ok := err.(*GameError); ok {
							log.Warn().
								Str("client_id", client.ID).
								Str("game_id", message.GameID).
								Float64("position", position).
								Str("error_code", gameErr.Code).
								Str("error_message", gameErr.Message).
								Msg("Invalid move attempt")

							response = &Message{
								Type:   "error",
								GameID: message.GameID,
								Error:  gameErr,
							}
							client.Conn.WriteJSON(response)
							continue
						}
					}

					log.Info().
						Str("client_id", client.ID).
						Str("game_id", message.GameID).
						Float64("position", position).
						Msg("Move successful")

					h.broadcastGameState(message.GameID)
				} else {
					log.Warn().
						Str("client_id", client.ID).
						Str("game_id", message.GameID).
						Interface("position", data["position"]).
						Msg("Invalid position format in move")

					response = &Message{
						Type:   "error",
						GameID: message.GameID,
						Error: &GameError{
							Code:    "INVALID_POSITION_FORMAT",
							Message: "Position must be a number",
						},
					}
					client.Conn.WriteJSON(response)
				}
			} else {
				log.Warn().
					Str("client_id", client.ID).
					Str("game_id", message.GameID).
					Interface("data", message.Data).
					Msg("Invalid move data format")

				response = &Message{
					Type:   "error",
					GameID: message.GameID,
					Error: &GameError{
						Code:    "INVALID_MOVE_FORMAT",
						Message: "Invalid move data format",
					},
				}
				client.Conn.WriteJSON(response)
			}
		default:
			log.Warn().
				Str("client_id", client.ID).
				Str("game_id", message.GameID).
				Str("message_type", message.Type).
				Msg("Unknown message type received")

			response = &Message{
				Type: "error",
				Error: &GameError{
					Code:    "UNKNOWN_MESSAGE_TYPE",
					Message: "Unknown message type",
				},
			}
			client.Conn.WriteJSON(response)
		}
	}
}

// broadcastGameState sends the current game state to all players in the game
func (h *Handler) broadcastGameState(gameID string) {
	h.manager.mutex.RLock()
	game, exists := h.manager.games[gameID]
	h.manager.mutex.RUnlock()

	if !exists {
		log.Warn().
			Str("game_id", gameID).
			Msg("Attempted to broadcast state for non-existent game")
		return
	}

	log.Debug().
		Str("game_id", gameID).
		Interface("game_state", game).
		Msg("Broadcasting game state")

	h.manager.broadcast <- &Message{
		Type:   "game_state",
		GameID: gameID,
		Data:   game,
	}
}
