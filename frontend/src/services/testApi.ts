// Test API service for debugging registration issues
import { User } from '../types';

// Store the JWT token in localStorage
export const setTestToken = (token: string): void => {
  localStorage.setItem('test_token', token);
};

export const getTestToken = (): string | null => {
  return localStorage.getItem('test_token');
};

// Use explicit backend URL
const API_URL = 'http://localhost:8082/api';

// Test registration function that bypasses the normal API flow
export const testRegister = async (username: string, password: string): Promise<User> => {
  console.log(`Attempting test registration for user: ${username}`);
  
  try {
    const response = await fetch(`${API_URL}/test/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      console.error(`Registration failed with status: ${response.status}`);
      const errorText = await response.text();
      console.error(`Error response: ${errorText}`);
      throw new Error(`Registration failed: ${response.statusText}`);
    }

    const data = await response.json();
    console.log('Registration successful:', data);
    return data;
  } catch (error) {
    console.error('Registration error:', error);
    throw error;
  }
};

// Test login function that bypasses the normal API flow
export const testLogin = async (username: string, password: string): Promise<{token: string; user: User}> => {
  console.log(`Attempting test login for user: ${username}`);
  
  try {
    // For testing purposes, we'll create a fake token and user
    // In a real implementation, this would make an API call to a test login endpoint
    const token = `test-token-${Math.random().toString(36).substring(2, 15)}`;
    const user = {
      id: `test-user-id-${Math.random().toString(36).substring(2, 15)}`,
      username: username
    };
    
    // Store the token
    setTestToken(token);
    
    console.log('Test login successful:', { token, user });
    return { token, user };
  } catch (error) {
    console.error('Test login error:', error);
    throw error;
  }
};
