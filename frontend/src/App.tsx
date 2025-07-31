import React, { useState } from 'react';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { ChatProvider } from './contexts/ChatContext';
import Login from './components/Login';
import Dashboard from './components/Dashboard';
import DebugPage from './pages/DebugPage';
import './App.css';

// Main application component that handles routing based on auth state
const AppContent: React.FC = () => {
  const { isAuthenticated, isLoading } = useAuth();
  const [showDebug, setShowDebug] = useState(false);

  // Toggle debug mode with Ctrl+Shift+D
  React.useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.ctrlKey && e.shiftKey && e.key === 'D') {
        setShowDebug(prev => !prev);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  // Show debug page if debug mode is enabled
  if (showDebug) {
    return <DebugPage />;
  }

  if (isLoading) {
    return (
      <div className="loading-container">
        <div className="loading-spinner"></div>
        <p>Loading...</p>
      </div>
    );
  }

  return isAuthenticated ? (
    <ChatProvider>
      <Dashboard />
    </ChatProvider>
  ) : (
    <Login />
  );
};

// Root component that provides context providers
function App() {
  return (
    <div className="App">
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </div>
  );
}

export default App;
