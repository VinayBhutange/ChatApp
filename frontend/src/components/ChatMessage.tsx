import React from 'react';
import { Message } from '../types';
import { generateAvatarColor } from '../services/avatarService';
import { useAuth } from '../contexts/AuthContext';
import '../styles/ChatMessage.css';

interface ChatMessageProps {
  message: Message;
}

const ChatMessage: React.FC<ChatMessageProps> = ({ message }) => {
  const { user } = useAuth();
  const isOwnMessage = user?.id === message.senderId;
  
  // Format timestamp
  const formatTime = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  return (
    <div className={`message-container ${isOwnMessage ? 'own-message' : 'other-message'}`}>
      {!isOwnMessage && (
        <div className="message-avatar" style={{ backgroundColor: generateAvatarColor(message.sender) }}>
          {message.sender.charAt(0).toUpperCase()}
        </div>
      )}
      <div className="message-bubble">
        {!isOwnMessage && (
          <div className="sender-name">{message.sender || 'Unknown User'}</div>
        )}
        <div className="message-content">{message.content}</div>
        <div className="message-time">{formatTime(message.timestamp)}</div>
      </div>
    </div>
  );
};

export default ChatMessage;
