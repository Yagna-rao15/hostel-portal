import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

const UpdatePassword = () => {
  const location = useLocation();
  const email = location.state?.email || ''; // Get email from route state
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setMessage(''); // Clear previous messages

    // Check if passwords match
    if (newPassword !== confirmPassword) {
      setMessage('Passwords do not match!');
      setLoading(false);
      return;
    }

    try {
      const response = await fetch('http://localhost:8080/update-password', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: email,
          password: newPassword,
        }),
      });

      const data = await response.json(); // Await this
      if (data.updated) {
        setMessage('Password updated successfully!');
        navigate('/login')
      } else {
        setMessage('Failed to update password: ' + data.message);
      }
    } catch (error) {
      console.error('Error updating password:', error);
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
          placeholder="Enter new password"
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          required
        />
        <br />
        <Input
          className="text-black"
          type="password" // Ensure this is a password field
          placeholder="Confirm new password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          required
        />
        <br />
        <Button
          className="w-96 h-1em text-base"
          type="submit"
          disabled={loading} // Disable button during loading
        >
          {loading ? 'Updating...' : 'Update Password'}
        </Button>
        {message && <p className="text-red-500">{message}</p>}
      </form>
    </div>
  );
}

export default UpdatePassword;
