import React, { useState, useEffect, useRef } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { generateAvatarColor } from '../services/avatarService';
import '../styles/IconSidebar.css';

const IconSidebar: React.FC = () => {
  const { user, logout } = useAuth();
  const [isMenuVisible, setIsMenuVisible] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsMenuVisible(false);
      }
    };

    if (isMenuVisible) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isMenuVisible]);

  return (
    <div className="icon-sidebar">
      <div className="sidebar-header">
        {/* Placeholder for a logo or brand icon */}
        <div className="logo">🚀</div>
      </div>
      <nav className="sidebar-nav">
        <a href="#" className="nav-item active" title="Chats">💬</a>
        <a href="#" className="nav-item" title="Contacts">👥</a>
        <a href="#" className="nav-item" title="Notifications">🔔</a>
        <a href="#" className="nav-item" title="Settings">⚙️</a>
      </nav>
      <div className="sidebar-footer" ref={menuRef}>
        {user && (
          <div 
            className="user-avatar-icon"
            style={{ backgroundColor: generateAvatarColor(user.username) }}
            onClick={() => setIsMenuVisible(!isMenuVisible)}
            title={user.username}
          >
            {user.username.charAt(0).toUpperCase()}
          </div>
        )}
        {isMenuVisible && user && (
          <div className="logout-menu">
            <div className="menu-user-info">
              <strong>{user.username}</strong>
            </div>
            <button className="logout-button-menu" onClick={logout}>
              Logout
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default IconSidebar;
