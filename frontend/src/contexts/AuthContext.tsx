import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { User, AuthState } from '../types';
import { getToken, setToken, clearToken, loginUser, registerUser } from '../services/api';
import { testRegister } from '../services/testApi';

interface AuthContextType extends AuthState {
  login: (username: string, password: string) => Promise<void>;
  register: (username: string, password: string) => Promise<void>;
  logout: () => void;
}

const initialState: AuthState = {
  user: null,
  token: null,
  isAuthenticated: false,
  isLoading: true,
  error: null,
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [state, setState] = useState<AuthState>(initialState);

  // Check for existing token on mount
  useEffect(() => {
    const checkAuth = async () => {
      const token = getToken();
      if (token) {
        try {
          // In a real app, you might want to validate the token with the server
          // For now, we'll just assume it's valid and set the user from localStorage if available
          const userJson = localStorage.getItem('user');
          if (userJson) {
            const user = JSON.parse(userJson);
            setState({
              user,
              token,
              isAuthenticated: true,
              isLoading: false,
              error: null,
            });
          } else {
            // If we have a token but no user, clear the token
            clearToken();
            setState({
              ...initialState,
              isLoading: false,
            });
          }
        } catch (error) {
          clearToken();
          setState({
            ...initialState,
            isLoading: false,
            error: 'Session expired. Please login again.',
          });
        }
      } else {
        setState({
          ...initialState,
          isLoading: false,
        });
      }
    };

    checkAuth();
  }, []);

  const login = async (username: string, password: string) => {
    setState({ ...state, isLoading: true, error: null });
    try {
      const response = await loginUser(username, password);
      localStorage.setItem('user', JSON.stringify(response.user));
      setState({
        user: response.user,
        token: response.token,
        isAuthenticated: true,
        isLoading: false,
        error: null,
      });
    } catch (error) {
      setState({
        ...state,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Failed to login',
      });
    }
  };

  const register = async (username: string, password: string) => {
    setState({ ...state, isLoading: true, error: null });
    try {
      console.log('Attempting registration with test endpoint');
      // Use the test registration endpoint instead of the regular one
      const user = await testRegister(username, password);
      console.log('Test registration successful, attempting login');
      // After registration, automatically log in
      await login(username, password);
    } catch (error) {
      console.error('Registration error:', error);
      setState({
        ...state,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Failed to register',
      });
    }
  };

  const logout = () => {
    clearToken();
    localStorage.removeItem('user');
    setState({
      ...initialState,
      isLoading: false,
    });
  };

  return (
    <AuthContext.Provider
      value={{
        ...state,
        login,
        register,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
