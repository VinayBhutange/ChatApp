import React from 'react';
import IconSidebar from './IconSidebar';
import ConversationList from './ConversationList';
import ChatWindow from './ChatWindow';
import '../styles/Dashboard.css';

const Dashboard: React.FC = () => {
  return (
    <div className="dashboard-layout">
      <IconSidebar />
      <ConversationList />
      <ChatWindow />
    </div>
  );
};

export default Dashboard;
