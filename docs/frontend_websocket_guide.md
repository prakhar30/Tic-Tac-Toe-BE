# Frontend WebSocket Integration Guide

This guide demonstrates how to integrate the Tic-Tac-Toe WebSocket APIs with your Next.js frontend application.

## Table of Contents
- [Setup](#setup)
- [WebSocket Connection Management](#websocket-connection-management)
- [Game State Management](#game-state-management)
- [Complete Example](#complete-example)
- [Error Handling](#error-handling)

## Setup

First, create a WebSocket service to manage the connection and message handling:

```typescript
// services/websocket.ts

interface WebSocketMessage {
  type: 'create_game' | 'join_game' | 'make_move' | 'game_state';
  gameId: string;
  data?: any;
  error?: GameError;
}

interface GameState {
  board: string[];
  players: { [key: string]: string };
  turn: string;
  winner: string;
  gameOver: boolean;
}

class WebSocketService {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectTimeout = 3000; // 3 seconds

  constructor(private token: string) {}

  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        this.ws = new WebSocket('ws://localhost:9092/ws');

        this.ws.onopen = () => {
          // Add authorization header
          const headers = {
            Authorization: `Bearer ${this.token}`,
          };
          
          // Send headers as part of the connection upgrade
          this.ws?.send(JSON.stringify(headers));
          
          console.log('WebSocket connected');
          this.reconnectAttempts = 0;
          resolve();
        };

        this.ws.onclose = () => {
          console.log('WebSocket disconnected');
          this.handleReconnect();
        };

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          reject(error);
        };
      } catch (error) {
        reject(error);
      }
    });
  }

  private handleReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);
      
      setTimeout(() => {
        this.connect().catch(console.error);
      }, this.reconnectTimeout);
    }
  }

  send(message: WebSocketMessage) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      console.error('WebSocket is not connected');
    }
  }

  onMessage(callback: (message: WebSocketMessage) => void) {
    if (this.ws) {
      this.ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        callback(message);
      };
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}

export default WebSocketService;
```

## Game State Management

Create a React context to manage the game state:

```typescript
// contexts/GameContext.tsx

import React, { createContext, useContext, useEffect, useState } from 'react';
import WebSocketService from '../services/websocket';

interface GameContextType {
  gameState: GameState | null;
  createGame: (gameId: string) => void;
  joinGame: (gameId: string) => void;
  makeMove: (position: number) => void;
  isConnected: boolean;
}

const GameContext = createContext<GameContextType | undefined>(undefined);

export function GameProvider({ children, token }: { children: React.ReactNode; token: string }) {
  const [ws, setWs] = useState<WebSocketService | null>(null);
  const [gameState, setGameState] = useState<GameState | null>(null);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    const websocket = new WebSocketService(token);
    
    websocket.connect()
      .then(() => {
        setIsConnected(true);
        setWs(websocket);
      })
      .catch(console.error);

    websocket.onMessage((message) => {
      if (message.type === 'game_state') {
        setGameState(message.data as GameState);
      }
    });

    return () => {
      websocket.disconnect();
    };
  }, [token]);

  const createGame = (gameId: string) => {
    ws?.send({
      type: 'create_game',
      gameId,
    });
  };

  const joinGame = (gameId: string) => {
    ws?.send({
      type: 'join_game',
      gameId,
    });
  };

  const makeMove = (position: number) => {
    if (!gameState) return;

    ws?.send({
      type: 'make_move',
      gameId: gameState.gameId,
      data: { position },
    });
  };

  return (
    <GameContext.Provider value={{ gameState, createGame, joinGame, makeMove, isConnected }}>
      {children}
    </GameContext.Provider>
  );
}

export const useGame = () => {
  const context = useContext(GameContext);
  if (!context) {
    throw new Error('useGame must be used within a GameProvider');
  }
  return context;
};
```

## Complete Example

Here's a complete example of a Tic-Tac-Toe game component:

```typescript
// components/TicTacToe.tsx

import { useGame } from '../contexts/GameContext';
import styles from './TicTacToe.module.css';

export default function TicTacToe() {
  const { gameState, makeMove, isConnected } = useGame();

  if (!isConnected) {
    return <div>Connecting to game server...</div>;
  }

  if (!gameState) {
    return <div>Waiting for game to start...</div>;
  }

  const handleCellClick = (position: number) => {
    if (gameState.board[position] === '' && !gameState.gameOver) {
      makeMove(position);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.status}>
        {gameState.gameOver 
          ? `Game Over! ${gameState.winner ? `Winner: ${gameState.winner}` : 'Draw!'}`
          : `Current Turn: ${gameState.turn}`
        }
      </div>
      
      <div className={styles.board}>
        {gameState.board.map((cell, index) => (
          <button
            key={index}
            className={styles.cell}
            onClick={() => handleCellClick(index)}
            disabled={cell !== '' || gameState.gameOver}
          >
            {cell}
          </button>
        ))}
      </div>

      <div className={styles.players}>
        {Object.entries(gameState.players).map(([playerId, symbol]) => (
          <div key={playerId}>
            Player {playerId}: {symbol}
          </div>
        ))}
      </div>
    </div>
  );
}
```

## Usage in Pages

Here's how to use the game components in your Next.js pages:

```typescript
// pages/game/[id].tsx

import { useEffect } from 'react';
import { useRouter } from 'next/router';
import { useGame } from '../../contexts/GameContext';
import TicTacToe from '../../components/TicTacToe';

export default function GamePage() {
  const router = useRouter();
  const { id: gameId } = router.query;
  const { joinGame } = useGame();

  useEffect(() => {
    if (gameId && typeof gameId === 'string') {
      joinGame(gameId);
    }
  }, [gameId]);

  return (
    <div>
      <h1>Game {gameId}</h1>
      <TicTacToe />
    </div>
  );
}
```

## Error Handling

Here are some common error scenarios and how to handle them:

1. Connection Errors:
```typescript
ws.onerror = (error) => {
  // Show error notification to user
  console.error('WebSocket error:', error);
  // Implement your error handling UI
};
```

2. Authentication Errors:
```typescript
ws.onclose = (event) => {
  if (event.code === 1008) { // Policy Violation (used for auth errors)
    // Handle authentication failure
    router.push('/login'); // Redirect to login page
  }
};
```

3. Game State Errors:
```typescript
// In your game component
useEffect(() => {
  if (gameState?.error) {
    // Show error message to user
    toast.error(gameState.error);
  }
}, [gameState]);
```

## Best Practices

1. **Connection Management**:
   - Always implement reconnection logic
   - Handle connection losses gracefully
   - Show appropriate UI feedback during connection states

2. **State Management**:
   - Keep WebSocket logic separate from UI components
   - Use React context for global game state
   - Implement proper cleanup in useEffect hooks

3. **Security**:
   - Always send the authentication token
   - Validate all incoming messages
   - Handle token expiration

4. **Error Handling**:
   - Implement comprehensive error handling
   - Show appropriate error messages to users
   - Log errors for debugging

5. **Performance**:
   - Clean up WebSocket connections when components unmount
   - Implement proper message queuing
   - Handle race conditions

## Testing

Here's an example of how to test your WebSocket integration:

```typescript
// __tests__/websocket.test.ts

import { render, act, fireEvent } from '@testing-library/react';
import { GameProvider } from '../contexts/GameContext';
import TicTacToe from '../components/TicTacToe';

describe('WebSocket Game Integration', () => {
  let mockWs: WebSocket;

  beforeEach(() => {
    mockWs = {
      send: jest.fn(),
      close: jest.fn(),
      addEventListener: jest.fn(),
      removeEventListener: jest.fn(),
    };
    
    // Mock WebSocket constructor
    global.WebSocket = jest.fn(() => mockWs);
  });

  it('should connect to WebSocket server', () => {
    render(
      <GameProvider token="test-token">
        <TicTacToe />
      </GameProvider>
    );

    expect(global.WebSocket).toHaveBeenCalledWith('ws://localhost:9092/ws');
  });

  // Add more tests...
});
``` 