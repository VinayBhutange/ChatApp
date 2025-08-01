import React, { useState } from 'react';
import { useChat } from '../contexts/ChatContext';
import { generateAvatarColor } from '../services/avatarService';
import '../styles/ConversationList.css';
import '../styles/Modal.css';

const ConversationList: React.FC = () => {
  const { rooms, joinRoom, currentRoom, createNewRoom } = useChat();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newRoomName, setNewRoomName] = useState('');

  const handleCreateRoom = async (e: React.FormEvent) => {
    e.preventDefault();
    if (newRoomName.trim()) {
      try {
        await createNewRoom(newRoomName.trim());
        setNewRoomName('');
        setIsModalOpen(false);
      } catch (error) {
        console.error('Failed to create room:', error);
        // Optionally, show an error message to the user in the modal
      }
    }
  };

  return (
    <div className="conversation-list">
      <div className="conversation-header">
        <h2>Chat</h2>
        {/* This button is currently for show; functionality can be added later */}
        <button className="new-message-btn" onClick={() => setIsModalOpen(true)}>+ New Room</button>
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

      {isModalOpen && (
        <div className="modal-backdrop">
          <div className="modal-content">
            <div className="modal-header">
              <h2>Create New Room</h2>
              <button className="modal-close-btn" onClick={() => setIsModalOpen(false)}>&times;</button>
            </div>
            <form onSubmit={handleCreateRoom}>
              <div className="modal-body">
                <div className="form-group">
                  <label htmlFor="roomName">Room Name</label>
                  <input
                    type="text"
                    id="roomName"
                    value={newRoomName}
                    onChange={(e) => setNewRoomName(e.target.value)}
                    placeholder="Enter a name for your new room"
                    autoFocus
                  />
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="modal-button secondary" onClick={() => setIsModalOpen(false)}>Cancel</button>
                <button type="submit" className="modal-button primary">Create</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default ConversationList;
