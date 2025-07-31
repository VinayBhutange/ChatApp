import React, { useState } from 'react';
import { testRegister, testLogin } from '../services/testApi';

const DebugRegister: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [response, setResponse] = useState<any>(null);
  const [loginResponse, setLoginResponse] = useState<any>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(null);
    setResponse(null);
    setLoginResponse(null);

    try {
      console.log(`Attempting debug registration for user: ${username}`);
      const result = await testRegister(username, password);
      console.log('Debug registration successful:', result);
      setSuccess(`User ${username} registered successfully!`);
      setResponse(result);
    } catch (err) {
      console.error('Debug registration error:', err);
      setError(err instanceof Error ? err.message : 'Failed to register');
    } finally {
      setLoading(false);
    }
  };
  
  const handleRegisterAndLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(null);
    setResponse(null);
    setLoginResponse(null);

    try {
      // First register the user
      console.log(`Attempting debug registration for user: ${username}`);
      const registerResult = await testRegister(username, password);
      console.log('Debug registration successful:', registerResult);
      setResponse(registerResult);
      
      // Then login with the same credentials
      console.log(`Attempting debug login for user: ${username}`);
      const loginResult = await testLogin(username, password);
      console.log('Debug login successful:', loginResult);
      setLoginResponse(loginResult);
      
      setSuccess(`User ${username} registered and logged in successfully!`);
    } catch (err) {
      console.error('Debug registration/login error:', err);
      setError(err instanceof Error ? err.message : 'Failed to register or login');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="debug-register">
      <h2>Debug Registration</h2>
      <form>
        <div>
          <label htmlFor="username">Username:</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="password">Password:</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <div className="button-group">
          <button type="button" onClick={handleSubmit} disabled={loading}>
            {loading ? 'Processing...' : 'Register Only'}
          </button>
          <button type="button" onClick={handleRegisterAndLogin} disabled={loading}>
            {loading ? 'Processing...' : 'Register & Login'}
          </button>
        </div>
      </form>

      {error && <div className="error">{error}</div>}
      {success && <div className="success">{success}</div>}
      
      {response && (
        <div className="response">
          <h3>Registration Response:</h3>
          <pre>{JSON.stringify(response, null, 2)}</pre>
        </div>
      )}
      
      {loginResponse && (
        <div className="response login-response">
          <h3>Login Response:</h3>
          <pre>{JSON.stringify(loginResponse, null, 2)}</pre>
        </div>
      )}
    </div>
  );
};

export default DebugRegister;
