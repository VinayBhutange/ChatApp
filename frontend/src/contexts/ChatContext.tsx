import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { Room, Message, ChatState } from '../types';
import { getRooms, createRoom, createWebSocketConnection } from '../services/api';
import { useAuth } from './AuthContext';

interface ChatContextType extends ChatState {
  fetchRooms: () => Promise<void>;
  createNewRoom: (name: string) => Promise<Room>;
  joinRoom: (roomId: string) => void;
  sendMessage: (content: string) => void;
  leaveRoom: () => void;
}

const initialState: ChatState = {
  rooms: [],
  currentRoom: null,
  messages: [],
  isLoading: false,
  error: null,
};

const ChatContext = createContext<ChatContextType | undefined>(undefined);

export const ChatProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [state, setState] = useState<ChatState>(initialState);
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const { isAuthenticated } = useAuth();

  // Fetch rooms when authenticated
  useEffect(() => {
    if (isAuthenticated) {
      fetchRooms();
    } else {
      setState(initialState);
    }
  }, [isAuthenticated]);

  // Clean up WebSocket connection when component unmounts
  useEffect(() => {
    return () => {
      if (socket) {
        socket.close();
      }
    };
  }, [socket]);

  // Handle WebSocket messages and reconnection
  useEffect(() => {
    if (socket && state.currentRoom) {
      const roomId = state.currentRoom.id;
      
      socket.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          setState((prevState) => ({
            ...prevState,
            messages: [...prevState.messages, message],
          }));
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };

      socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        setState((prevState) => ({
          ...prevState,
          error: 'WebSocket connection error - will attempt to reconnect',
        }));
      };

      socket.onclose = (event) => {
        console.log(`WebSocket connection closed: ${event.code} ${event.reason || ''}`);
        
        // Attempt to reconnect after a delay if not a clean close
        if (!event.wasClean && state.currentRoom) {
          console.log(`Attempting to reconnect to room ${roomId} in 3 seconds...`);
          
          // Set a timeout to reconnect
          const reconnectTimer = setTimeout(() => {
            console.log(`Reconnecting to room ${roomId}...`);
            try {
              const newSocket = createWebSocketConnection(roomId);
              setSocket(newSocket);
              console.log('Reconnection attempt initiated');
            } catch (reconnectError) {
              console.error('Failed to reconnect:', reconnectError);
            }
          }, 3000);
          
          // Clear the timeout if the component unmounts or socket changes
          return () => clearTimeout(reconnectTimer);
        }
      };
    }
  }, [socket, state.currentRoom]);

  const fetchRooms = async () => {
    setState((prevState) => ({ ...prevState, isLoading: true, error: null }));
    try {
      const response = await getRooms();
      const rooms = response.rooms || [];
      setState((prevState) => ({
        ...prevState,
        rooms,
        isLoading: false,
      }));
    } catch (error) {
      setState((prevState) => ({
        ...prevState,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Failed to fetch rooms',
      }));
    }
  };

  const createNewRoom = async (name: string) => {
    setState((prevState) => ({ ...prevState, isLoading: true, error: null }));
    try {
      const newRoom = await createRoom(name);
      setState((prevState) => ({
        ...prevState,
        rooms: [...prevState.rooms, newRoom],
        isLoading: false,
      }));
      return newRoom;
    } catch (error) {
      setState((prevState) => ({
        ...prevState,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Failed to create room',
      }));
      throw error;
    }
  };

  const joinRoom = (roomId: string) => {
    // Close existing socket if any
    if (socket) {
      socket.close();
    }

    const room = state.rooms.find((r) => r.id === roomId);
    if (!room) {
      setState((prevState) => ({
        ...prevState,
        error: 'Room not found',
      }));
      return;
    }

    try {
      // Create new WebSocket connection
      const newSocket = createWebSocketConnection(roomId);
      setSocket(newSocket);

      // Update state
      setState((prevState) => ({
        ...prevState,
        currentRoom: room,
        messages: [], // Clear previous messages
        error: null,
      }));
    } catch (error) {
      setState((prevState) => ({
        ...prevState,
        error: error instanceof Error ? error.message : 'Failed to join room',
      }));
    }
  };

  const sendMessage = (content: string) => {
    if (!socket || socket.readyState !== WebSocket.OPEN || !state.currentRoom) {
      setState((prevState) => ({
        ...prevState,
        error: 'Not connected to a room',
      }));
      return;
    }

    const message = {
      content,
    };

    socket.send(JSON.stringify(message));
  };

  const leaveRoom = () => {
    if (socket) {
      socket.close();
      setSocket(null);
    }

    setState((prevState) => ({
      ...prevState,
      currentRoom: null,
      messages: [],
    }));
  };

  return (
    <ChatContext.Provider
      value={{
        ...state,
        fetchRooms,
        createNewRoom,
        joinRoom,
        sendMessage,
        leaveRoom,
      }}
    >
      {children}
    </ChatContext.Provider>
  );
};

export const useChat = (): ChatContextType => {
  const context = useContext(ChatContext);
  if (context === undefined) {
    throw new Error('useChat must be used within a ChatProvider');
  }
  return context;
};
