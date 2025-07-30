// API service for handling backend communication
import { User, Room, Message } from '../types';

// Using relative URL with proxy configuration from package.json
const API_URL = '/api';

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
  
  // Create WebSocket URL using window.location to match the current host
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const host = window.location.host;
  // Include token directly in the URL for authentication
  const ws = new WebSocket(`${protocol}//${host}/api/ws?room_id=${roomId}&token=${token}`);
  
  // Set up connection event handlers
  ws.onopen = () => {
    console.log('WebSocket connection established');
  };
  
  return ws;
};
