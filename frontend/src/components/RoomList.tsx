import React, { useState, useEffect } from 'react';
import { useChat } from '../contexts/ChatContext';
import '../styles/RoomList.css';

const RoomList: React.FC = () => {
  const { rooms, fetchRooms, createNewRoom, joinRoom, currentRoom, isLoading, error } = useChat();
  const [newRoomName, setNewRoomName] = useState('');
  const [showCreateForm, setShowCreateForm] = useState(false);

  useEffect(() => {
    fetchRooms();
  }, []);

  const handleCreateRoom = async (e: React.FormEvent) => {
    e.preventDefault();
    if (newRoomName.trim()) {
      try {
        await createNewRoom(newRoomName.trim());
        setNewRoomName('');
        setShowCreateForm(false);
      } catch (error) {
        console.error('Failed to create room:', error);
      }
    }
  };

  const handleJoinRoom = (roomId: string) => {
    joinRoom(roomId);
  };

  return (
    <div className="room-list-container">
      <div className="room-list-header">
        <h2>Chat Rooms</h2>
        <button 
          className="create-room-button"
          onClick={() => setShowCreateForm(!showCreateForm)}
        >
          {showCreateForm ? 'Cancel' : '+ New Room'}
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}

      {showCreateForm && (
        <form className="create-room-form" onSubmit={handleCreateRoom}>
          <input
            type="text"
            placeholder="Room name"
            value={newRoomName}
            onChange={(e) => setNewRoomName(e.target.value)}
            required
          />
          <button type="submit" disabled={isLoading}>
            {isLoading ? 'Creating...' : 'Create'}
          </button>
        </form>
      )}

      {isLoading && !showCreateForm ? (
        <div className="loading">Loading rooms...</div>
      ) : (
        <ul className="room-list">
          {rooms.length === 0 ? (
            <li className="no-rooms">No rooms available. Create one!</li>
          ) : (
            rooms.map((room) => (
              <li 
                key={room.id} 
                className={`room-item ${currentRoom?.id === room.id ? 'active' : ''}`}
                onClick={() => handleJoinRoom(room.id)}
              >
                {room.name}
              </li>
            ))
          )}
        </ul>
      )}
    </div>
  );
};

export default RoomList;
