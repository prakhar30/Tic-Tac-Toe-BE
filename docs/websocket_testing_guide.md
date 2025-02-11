# Testing WebSocket Connections with Postman

This guide demonstrates how to test the Tic-Tac-Toe WebSocket API using Postman.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Setting Up Postman](#setting-up-postman)
- [Authentication](#authentication)
- [WebSocket Connection](#websocket-connection)
- [Testing Game Actions](#testing-game-actions)
- [Common Test Scenarios](#common-test-scenarios)
- [Troubleshooting](#troubleshooting)

## Prerequisites

1. Install Postman (latest version)
2. Have the backend server running locally
3. Have a valid authentication token (obtained through the login API)

## Setting Up Postman

1. Create a new WebSocket request:
   - Click "New" → "WebSocket Request"
   - Enter the WebSocket URL: `ws://localhost:9092/ws`

2. Set up environment variables:
   ```
   WEBSOCKET_URL: ws://localhost:9092/ws
   AUTH_TOKEN: <your_auth_token>
   ```

## Authentication

Before testing WebSocket connections, you need to obtain an authentication token:

1. Create a login request:
   ```http
   POST http://localhost:9091/v1/users/login
   Content-Type: application/json

   {
     "username": "your_username",
     "password": "your_password"
   }
   ```

2. Save the returned token in your Postman environment:
   ```javascript
   // In the "Tests" tab of the login request
   pm.environment.set("AUTH_TOKEN", pm.response.json().access_token);
   ```

## WebSocket Connection

1. Configure WebSocket headers:
   ```
   Authorization: Bearer {{AUTH_TOKEN}}
   ```

2. Connect to WebSocket:
   - Click "Connect"
   - You should see a "Connected" status

## Testing Game Actions

### 1. Create a New Game

```json
{
  "type": "create_game",
  "gameId": "test_game_123",
  "playerId": "player1"
}
```

Expected Response:
```json
{
  "type": "game_state",
  "gameId": "test_game_123",
  "data": {
    "board": ["","","","","","","","",""],
    "players": {
      "player1": "X"
    },
    "turn": "player1",
    "winner": "",
    "gameOver": false
  }
}
```

### 2. Join a Game

```json
{
  "type": "join_game",
  "gameId": "test_game_123",
  "playerId": "player2"
}
```

Expected Response:
```json
{
  "type": "game_state",
  "gameId": "test_game_123",
  "data": {
    "board": ["","","","","","","","",""],
    "players": {
      "player1": "X",
      "player2": "O"
    },
    "turn": "player1",
    "winner": "",
    "gameOver": false
  }
}
```

### 3. Make a Move

```json
{
  "type": "make_move",
  "gameId": "test_game_123",
  "playerId": "player1",
  "data": {
    "position": 4
  }
}
```

Expected Response:
```json
{
  "type": "game_state",
  "gameId": "test_game_123",
  "data": {
    "board": ["","","","","X","","","",""],
    "players": {
      "player1": "X",
      "player2": "O"
    },
    "turn": "player2",
    "winner": "",
    "gameOver": false
  }
}
```

## Common Test Scenarios

### 1. Testing Authentication

1. Connect without token:
   - Expected: Connection refused with 401 Unauthorized

2. Connect with invalid token:
   - Expected: Connection refused with 401 Unauthorized

3. Connect with expired token:
   - Expected: Connection refused with 401 Unauthorized

### 2. Testing Game Rules

1. Making a move out of turn:
```json
{
  "type": "make_move",
  "gameId": "test_game_123",
  "playerId": "player2",  // When it's player1's turn
  "data": {
    "position": 0
  }
}
```
Expected: Move rejected, no state change

2. Making a move in an occupied position:
```json
{
  "type": "make_move",
  "gameId": "test_game_123",
  "playerId": "player2",
  "data": {
    "position": 4  // Already occupied
  }
}
```
Expected: Move rejected, no state change

3. Joining a full game:
```json
{
  "type": "join_game",
  "gameId": "test_game_123",
  "playerId": "player3"  // Game already has 2 players
}
```
Expected: Join rejected, no state change

### 3. Testing Win Conditions

1. Horizontal win:
```json
// Sequence of moves to test horizontal win
[
  {"position": 0, "playerId": "player1"},
  {"position": 3, "playerId": "player2"},
  {"position": 1, "playerId": "player1"},
  {"position": 4, "playerId": "player2"},
  {"position": 2, "playerId": "player1"}
]
```
Expected: Game over with player1 as winner

2. Vertical win:
```json
// Sequence of moves to test vertical win
[
  {"position": 0, "playerId": "player1"},
  {"position": 1, "playerId": "player2"},
  {"position": 3, "playerId": "player1"},
  {"position": 4, "playerId": "player2"},
  {"position": 6, "playerId": "player1"}
]
```
Expected: Game over with player1 as winner

3. Diagonal win:
```json
// Sequence of moves to test diagonal win
[
  {"position": 0, "playerId": "player1"},
  {"position": 1, "playerId": "player2"},
  {"position": 4, "playerId": "player1"},
  {"position": 3, "playerId": "player2"},
  {"position": 8, "playerId": "player1"}
]
```
Expected: Game over with player1 as winner

## Troubleshooting

### Common Issues and Solutions

1. Connection Refused
   - Check if the server is running
   - Verify the WebSocket URL
   - Ensure the port (9092) is not blocked

2. Authentication Failed
   - Check if the token is valid
   - Verify the token format in the Authorization header
   - Check token expiration

3. Messages Not Received
   - Check WebSocket connection status
   - Verify message format
   - Check console for error messages

### Debugging Tips

1. Use Postman Console:
   - View → Show Postman Console
   - Monitor WebSocket messages and errors

2. Enable Verbose Logging:
   ```javascript
   // In the "Pre-request Script" tab
   console.log('Sending message:', pm.request.body);
   ```

3. Test Connection Health:
   ```javascript
   // In the "Tests" tab
   pm.test("WebSocket is connected", () => {
     pm.expect(pm.response.code).to.equal(101);
   });
   ```

## Best Practices

1. Always clean up after testing:
   - Close WebSocket connections
   - Remove test game data

2. Test with multiple clients:
   - Open multiple Postman tabs
   - Simulate real multiplayer scenarios

3. Validate responses:
   - Check message types
   - Verify game state changes
   - Confirm turn order

4. Test error cases:
   - Invalid moves
   - Malformed messages
   - Connection interruptions 