// API service for handling backend communication
import { User, Room, Message } from '../types';

// Use explicit backend URL instead of relying on proxy
const API_URL = 'http://localhost:8082/api';

// Store the JWT token in localStorage
export const setToken = (token: string): void => {
  localStorage.setItem('token', token);
};

export const getToken = (): string | null => {
  return localStorage.getItem('token');
};

export const clearToken = (): void => {
  localStorage.removeItem('token');
};

// Helper function for API requests with retry logic
const apiRequest = async (
  endpoint: string,
  method: string = 'GET',
  body?: any,
  requiresAuth: boolean = true,
  retries: number = 3
): Promise<any> => {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  };

  if (requiresAuth) {
    const token = getToken();
    if (!token) {
      throw new Error('Authentication required');
    }
    headers['Authorization'] = `Bearer ${token}`;
  }

  const options: RequestInit = {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
    // Add mode and credentials for CORS support
    mode: 'cors',
    credentials: 'include'
  };

  let lastError: Error | null = null;
  
  for (let attempt = 0; attempt < retries; attempt++) {
    try {
      console.log(`Attempting API request to ${API_URL}${endpoint} (attempt ${attempt + 1}/${retries})`);
      const response = await fetch(`${API_URL}${endpoint}`, options);

      if (!response.ok) {
        const errorText = await response.text();
        let errorMessage = response.statusText;
        try {
          const errorData = JSON.parse(errorText);
          errorMessage = errorData.message || errorMessage;
        } catch (e) {
          // If the error response is not valid JSON, use the text as is
          errorMessage = errorText || errorMessage;
        }
        throw new Error(errorMessage);
      }
      
      return await response.json();
    } catch (error) {
      lastError = error instanceof Error ? error : new Error(String(error));
      console.error(`API request failed (attempt ${attempt + 1}/${retries}):`, lastError);
      
      if (attempt < retries - 1) {
        // Wait before retrying (exponential backoff)
        const delay = Math.pow(2, attempt) * 500;
        await new Promise(resolve => setTimeout(resolve, delay));
      }
    }
  }
  
  // If we get here, all retries failed
  throw lastError || new Error('API request failed after multiple attempts');
};

// Auth API calls
export const registerUser = async (username: string, password: string): Promise<User> => {
  return apiRequest('/register', 'POST', { username, password }, false);
};

export const loginUser = async (username: string, password: string): Promise<{ token: string; user: User }> => {
  const response = await apiRequest('/login', 'POST', { username, password }, false);
  setToken(response.token);
  return response;
};

export const logout = (): void => {
  clearToken();
};

// Room API calls
export const getRooms = async (): Promise<Room[]> => {
  const response = await apiRequest('/rooms');
  return response.rooms;
};

export const createRoom = async (name: string): Promise<Room> => {
  return apiRequest('/rooms/create', 'POST', { name });
};

// WebSocket connection for chat with retry logic
export const createWebSocketConnection = (roomId: string): WebSocket => {
  const token = getToken();
  if (!token) {
    throw new Error('Authentication required for WebSocket connection');
  }
  
  // Use explicit WebSocket URL
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  // Use 127.0.0.1 instead of localhost to avoid potential DNS issues
  const host = '127.0.0.1:8082'; // Hardcode the backend server address
  
  // Include token directly in the URL for authentication
  const wsUrl = `${protocol}//${host}/api/ws?room_id=${roomId}&token=${token}`;
  console.log(`Attempting to connect to WebSocket at: ${wsUrl}`);
  console.log(`Token length: ${token.length}, Room ID: ${roomId}`);
  
  // Log network status and browser information
  console.log(`Network status - Online: ${navigator.onLine}`);
  console.log(`Browser: ${navigator.userAgent}`);
  
  // First try a simple fetch to check if the server is reachable
  fetch(`http://${host}/api/rooms`)
    .then(response => {
      console.log('Backend server is reachable via HTTP:', response.status);
    })
    .catch(err => {
      console.error('Backend server is not reachable via HTTP:', err);
    });
  
  try {
    // Create WebSocket with custom headers if needed
    const ws = new WebSocket(wsUrl);
    
    // Set up connection event handlers
    ws.onopen = (event) => {
      console.log('WebSocket connection established successfully', event);
    };
    
    ws.onerror = (error) => {
      console.error('WebSocket connection error:', error);
      // Try to diagnose the error
      if (!navigator.onLine) {
        console.error('Network appears to be offline');
      }
      
      // Check if the backend server is reachable
      fetch(`http://${host}/api/rooms`)
        .then(response => {
          console.log('Backend server is reachable:', response.status);
        })
        .catch(err => {
          console.error('Backend server is not reachable:', err);
        });
    };
    
    // Set up reconnection logic
    ws.onclose = (event) => {
      console.log(`WebSocket connection closed: Code ${event.code} Reason: "${event.reason || 'No reason provided'}" Clean: ${event.wasClean}`);
      
      // Attempt to reconnect after a delay, unless it was a clean close
      if (!event.wasClean) {
        const reconnectDelay = 3000; // 3 seconds
        console.log(`Connection closed unexpectedly, attempting to reconnect in ${reconnectDelay/1000} seconds...`);
        
        setTimeout(() => {
          try {
            console.log('Attempting to reconnect WebSocket...');
            // Note: We're not actually creating a new connection here to avoid infinite recursion
            // The ChatContext component should handle reconnection by calling this function again
          } catch (error) {
            console.error('Failed to reconnect WebSocket:', error);
          }
        }, reconnectDelay);
      }
    };
    
    // Add a message handler for debugging
    ws.onmessage = (event) => {
      console.log('WebSocket message received:', event.data);
      // The actual message handling will be done by the ChatContext
    };
    
    return ws;
  } catch (error) {
    console.error('Failed to create WebSocket connection:', error);
    throw error;
  }
};
