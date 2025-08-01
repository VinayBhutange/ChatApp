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
        // Optionally, display an error message to the user
      }
    }
  };

  return (
    <div className="conversation-list">
      <div className="conversation-list-header">
        <h2>Conversations</h2>
        <button className="new-room-btn" onClick={() => setIsModalOpen(true)}>+</button>
      </div>
      <ul>
        {rooms.map((room) => (
          <li
            key={room.id}
            className={`conversation-item ${currentRoom?.id === room.id ? 'active' : ''}`}
            onClick={() => joinRoom(room.id)}
          >
            <div className="avatar" style={{ backgroundColor: generateAvatarColor(room.name) }}>
              {room.name.charAt(0).toUpperCase()}
            </div>
            <div className="conversation-details">
              <span className="conversation-name">{room.name}</span>
            </div>
          </li>
        ))}
      </ul>

      {isModalOpen && (
        <div className="modal-overlay">
          <div className="modal">
            <div className="modal-header">
              <h2>Create New Room</h2>
              <button className="modal-close-btn" onClick={() => setIsModalOpen(false)}>&times;</button>
            </div>
            <div className="modal-body">
              <form onSubmit={handleCreateRoom}>
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
                <div className="modal-footer">
                  <button type="button" className="modal-button secondary" onClick={() => setIsModalOpen(false)}>Cancel</button>
                  <button type="submit" className="modal-button primary">Create</button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default ConversationList;
