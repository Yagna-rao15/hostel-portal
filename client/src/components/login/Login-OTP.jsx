import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button"
import { InputOTP, InputOTPGroup, InputOTPSlot, InputOTPSeparator, } from "@/components/ui/input-otp";

const Login_Email = () => {
  const [email, setEmail] = useState('');
  const [otp, setOtp] = useState('');
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const response = await fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: email,
          password: password
        })
      });

      const data = await response.json(); // Await this
      if (!data.valid) {
        setMessage('Enter your valid institute mail address');
      } else if (data.isRegistered) {
        navigate('/login/password');
      } else {
        navigate('/login/verify-email');
      }

    } catch (error) {
      console.error('Error checking email:', error);
      setMessage('Server error. Please try again later.');
    }
    setLoading(false); // Set loading to false after check completes
  };

  const handleOtpChange = (value) => {
    setOtp(value);
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && otp.length === 6) {
      e.preventDefault();
      navigate('/login/password');
      formRef.current.submit();
    }
  };

  return (
    <div>
      <form onSubmit={handleSubmit} className="flex flex-col items-center space-y-4">
        <Input
          className="text-black"
          type="email"
          placeholder="Type your email address"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
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
    </div>
  );
}

export default Login_Email;


