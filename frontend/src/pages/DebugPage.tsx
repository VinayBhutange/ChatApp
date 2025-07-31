import React from 'react';
import DebugRegister from '../components/DebugRegister';
import '../styles/Debug.css';

const DebugPage: React.FC = () => {
  return (
    <div className="debug-page">
      <h1>Debug Tools</h1>
      <p>Use these tools to test and debug the application functionality.</p>
      
      <DebugRegister />
    </div>
  );
};

export default DebugPage;
