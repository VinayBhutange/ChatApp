import React from 'react';
import { useChat } from '../contexts/ChatContext';
import { generateAvatarColor } from '../services/avatarService';
import '../styles/ConversationList.css';

const ConversationList: React.FC = () => {
  const { rooms, joinRoom, currentRoom } = useChat();

  return (
    <div className="conversation-list">
      <div className="conversation-header">
        <h2>Chat</h2>
        {/* This button is currently for show; functionality can be added later */}
        <button className="new-message-btn">+ New Room</button>
      </div>
      <div className="search-bar">
        <input type="text" placeholder="Search" />
      </div>
      <div className="conversations">
        {rooms.map(room => (
          <div 
            key={room.id} 
            className={`conversation-item ${currentRoom?.id === room.id ? 'active' : ''}`}
            onClick={() => joinRoom(room.id)}
          >
            <div className="avatar" style={{ backgroundColor: generateAvatarColor(room.name) }}>
              {room.name.charAt(0).toUpperCase()}
            </div>
            <div className="conversation-details">
              <div className="conversation-name-time">
                <span className="name">{room.name}</span>
              </div>
              <div className="conversation-message">
                <p>Click to join this room</p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ConversationList;
