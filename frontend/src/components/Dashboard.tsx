import React from 'react';
import { useAuth } from '../contexts/AuthContext';
import RoomList from './RoomList';
import ChatRoom from './ChatRoom';
import '../styles/Dashboard.css';

const Dashboard: React.FC = () => {
  const { user, logout } = useAuth();

  return (
    <div className="dashboard-container">
      <header className="dashboard-header">
        <h1>Chat App</h1>
        <div className="user-info">
          <span>Welcome, {user?.username}</span>
          <button className="logout-button" onClick={logout}>
            Logout
          </button>
        </div>
      </header>
      
      <main className="dashboard-content">
        <aside className="sidebar">
          <RoomList />
        </aside>
        
        <section className="main-content">
          <ChatRoom />
        </section>
      </main>
    </div>
  );
};

export default Dashboard;
