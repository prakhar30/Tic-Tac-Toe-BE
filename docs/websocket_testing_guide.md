# WebSocket Testing Guide

This guide provides examples of how to test the WebSocket functionality of the Tic-Tac-Toe game.

## Prerequisites

1. The server is running on `ws://localhost:9092/ws`
2. You have valid authentication tokens for testing (one for each player)

## Connection Setup

1. Connect to the WebSocket endpoint with an authentication token:
```javascript
const ws = new WebSocket('ws://localhost:9092/ws');
ws.onopen = () => {
  // Send authentication token in headers
  ws.send(JSON.stringify({
    Authorization: `Bearer ${token}`
  }));
};
```

## Test Scenarios

### 1. Creating a Game

Request:
```json
{
  "type": "create_game",
  "gameId": "test_game_123"
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
      "alice": "X"  // Username from token
    },
    "turn": "alice",
    "winner": "",
    "gameOver": false,
    "gameReady": false
  }
}
```

### 2. Joining a Game

Request:
```json
{
  "type": "join_game",
  "gameId": "test_game_123"
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
      "alice": "X",
      "bob": "O"  // Second player's username from token
    },
    "turn": "alice",
    "winner": "",
    "gameOver": false,
    "gameReady": true
  }
}
```

### 3. Making a Move

Request:
```json
{
  "type": "make_move",
  "gameId": "test_game_123",
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
      "alice": "X",
      "bob": "O"
    },
    "turn": "bob",
    "winner": "",
    "gameOver": false,
    "gameReady": true
  }
}
```

## Testing Error Cases

### 1. Moving Out of Turn
```json
{
  "type": "make_move",
  "gameId": "test_game_123",
  "data": {
    "position": 0
  }
}
```
Expected Error: "NOT_PLAYERS_TURN"

### 2. Moving to Occupied Position
```json
{
  "type": "make_move",
  "gameId": "test_game_123",
  "data": {
    "position": 4
  }
}
```
Expected Error: "POSITION_OCCUPIED"

### 3. Joining Full Game
```json
{
  "type": "join_game",
  "gameId": "test_game_123"
}
```
Expected Error: "GAME_FULL"

## Testing Win Conditions

### 1. Horizontal Win
Sequence of moves:
```json
[
  {"position": 0},  // X
  {"position": 3},  // O
  {"position": 1},  // X
  {"position": 4},  // O
  {"position": 2}   // X wins
]
```

### 2. Vertical Win
Sequence of moves:
```json
[
  {"position": 0},  // X
  {"position": 1},  // O
  {"position": 3},  // X
  {"position": 4},  // O
  {"position": 6}   // X wins
]
```

### 3. Diagonal Win
Sequence of moves:
```json
[
  {"position": 0},  // X
  {"position": 1},  // O
  {"position": 4},  // X
  {"position": 3},  // O
  {"position": 8}   // X wins
]
```

## Cleanup

1. Close WebSocket connections after testing:
```javascript
ws.close();
```

2. Verify disconnection events are handled properly.

## Tips

1. Use different browser tabs or Postman to test multiplayer scenarios
2. Keep authentication tokens ready for different test users
3. Monitor the server logs for detailed information about game state changes
4. Test reconnection scenarios by temporarily disconnecting clients
5. Verify error messages are properly displayed to users
6. Test with invalid game IDs and positions to ensure proper error handling 