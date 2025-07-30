// User type definition
export interface User {
  id: string;
  username: string;
}

// Chat room type definition
export interface Room {
  id: string;
  name: string;
}

// Message type definition
export interface Message {
  id: string;
  roomId: string;
  senderId: string;
  content: string;
  timestamp: string;
  senderUsername?: string; // Optional field that might be populated from the API
}

// Authentication state
export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

// Chat state
export interface ChatState {
  rooms: Room[];
  currentRoom: Room | null;
  messages: Message[];
  isLoading: boolean;
  error: string | null;
}
