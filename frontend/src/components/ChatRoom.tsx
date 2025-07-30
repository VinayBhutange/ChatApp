import React, { useState, useRef, useEffect } from 'react';
import { useChat } from '../contexts/ChatContext';
import ChatMessage from './ChatMessage';
import '../styles/ChatRoom.css';

const ChatRoom: React.FC = () => {
  const { currentRoom, messages, sendMessage, leaveRoom, error } = useChat();
  const [newMessage, setNewMessage] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Scroll to bottom when new messages arrive
  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (newMessage.trim()) {
      sendMessage(newMessage.trim());
      setNewMessage('');
    }
  };

  if (!currentRoom) {
    return (
      <div className="chat-room-placeholder">
        <h2>Select a room to start chatting</h2>
      </div>
    );
  }

  return (
    <div className="chat-room-container">
      <div className="chat-room-header">
        <h2>{currentRoom.name}</h2>
        <button className="leave-room-button" onClick={leaveRoom}>
          Leave Room
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}

      <div className="messages-container">
        {messages.length === 0 ? (
          <div className="no-messages">No messages yet. Be the first to say hello!</div>
        ) : (
          messages.map((message) => (
            <ChatMessage key={message.id} message={message} />
          ))
        )}
        <div ref={messagesEndRef} />
      </div>

      <form className="message-form" onSubmit={handleSendMessage}>
        <input
          type="text"
          placeholder="Type a message..."
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
        />
        <button type="submit">Send</button>
      </form>
    </div>
  );
};

export default ChatRoom;
