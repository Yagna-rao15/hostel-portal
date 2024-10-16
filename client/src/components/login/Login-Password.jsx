import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

const Login_Email = () => {
  const location = useLocation();
  const email = location.state?.email || ''; // Get email from route state
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const response = await fetch('http://localhost:8080/password', { // Make sure this matches your server endpoint
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ // Send email and password
          email: email,
          password: password,
        })
      });

      const data = await response.json(); // Await this
      if (!data.valid) {
        setMessage('Wrong Email or Password Combination');
      } else {
        navigate('/form', { state: { email } }); // Navigate to the correct route
      }

    } catch (error) {
      console.error('Error checking email:', error);
      setMessage('Server error. Please try again later.');
    }
    setLoading(false); // Set loading to false after check completes
  };

  return (
    <div>
      <form onSubmit={handleSubmit}>
        <Input
          className="text-black"
          type="password" // Ensure this is a password field
          placeholder="Type your password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <br />
        <Button
          className="w-96 h-1em text-base"
          type="submit"
          disabled={loading} // Disable button during loading
        >
          {loading ? 'Verifying...' : 'Enter'}
        </Button>
        {message && <p className="text-red-500">{message}</p>}
      </form>
      <div className="text-right mb-3 mt-3">
        <a href="/" onClick={(e) => {
          e.preventDefault();
          navigate('/login/update-password', { state: { email } });
        }} className="text-gray-800">Forgot password?</a>
      </div>
    </div>
  );
}

export default Login_Email;

