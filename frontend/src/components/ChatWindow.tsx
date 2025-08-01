import React, { useState, useEffect, useRef } from 'react';
import { useChat } from '../contexts/ChatContext';
import { generateAvatarColor } from '../services/avatarService';
import '../styles/ChatWindow.css';
import ChatMessage from './ChatMessage';

const ChatWindow: React.FC = () => {
  const { currentRoom, messages, sendMessage } = useChat();
  const [newMessage, setNewMessage] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (newMessage.trim() && currentRoom) {
      sendMessage(newMessage.trim());
      setNewMessage('');
    }
  };

  if (!currentRoom) {
    return (
      <div className="chat-window-placeholder">
        <h2>Select a conversation to start chatting</h2>
      </div>
    );
  }

  return (
    <div className="chat-window">
      <div className="chat-header">
        <div className="contact-info">
          <div className="avatar online" style={{ backgroundColor: generateAvatarColor(currentRoom.name) }}>{currentRoom.name.charAt(0).toUpperCase()}</div>
          <div className="contact-details">
            <span className="name">{currentRoom.name}</span>
            <span className="status">Online</span>
          </div>
        </div>
        <div className="header-actions">
          <button>⭐</button>
          <button>︙</button>
        </div>
      </div>
      <div className="message-area">
        {messages.map(msg => (
          <ChatMessage key={msg.id} message={msg} />
        ))}
        <div ref={messagesEndRef} />
      </div>
      <form className="message-input-area" onSubmit={handleSendMessage}>
        <input 
          type="text" 
          placeholder="Write your message" 
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
        />
        <button type="submit" className="send-button">➢</button>
      </form>
    </div>
  );
};

export default ChatWindow;
