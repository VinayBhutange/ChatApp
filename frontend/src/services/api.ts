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

// Helper function for API requests
const apiRequest = async (
  endpoint: string,
  method: string = 'GET',
  body?: any,
  requiresAuth: boolean = true
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
  };

  const response = await fetch(`${API_URL}${endpoint}`, options);

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.message || response.statusText);
  }

  return response.json();
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

// WebSocket connection for chat
export const createWebSocketConnection = (roomId: string): WebSocket => {
  const token = getToken();
  if (!token) {
    throw new Error('Authentication required for WebSocket connection');
  }
  
  // Use explicit WebSocket URL
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const host = 'localhost:8082'; // Hardcode the backend server address
  
  // Include token directly in the URL for authentication
  const wsUrl = `${protocol}//${host}/api/ws?room_id=${roomId}&token=${token}`;
  console.log(`Attempting to connect to WebSocket at: ${wsUrl}`);
  
  try {
    const ws = new WebSocket(wsUrl);
    
    // Set up connection event handlers
    ws.onopen = () => {
      console.log('WebSocket connection established successfully');
    };
    
    ws.onerror = (error) => {
      console.error('WebSocket connection error:', error);
    };
    
    return ws;
  } catch (error) {
    console.error('Failed to create WebSocket connection:', error);
    throw error;
  }
};
