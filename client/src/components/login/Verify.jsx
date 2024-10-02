import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { REGEXP_ONLY_DIGITS_AND_CHARS } from "input-otp";
import { InputOTP, InputOTPGroup, InputOTPSlot, InputOTPSeparator, } from "@/components/ui/input-otp";

const VerifyEmailPage = () => {
  const [otp, setOtp] = useState('');
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const response = await fetch('http://localhost:8080/api/verify-email', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ otp })
      });

      const data = await response.json();
      if (!data.valid) {
        setMessage('Enter your valid institute mail address');
      } else if (data.isRegistered) {
        navigate('/login/password');
      } else {
        navigate('/login/verify-email');
      }
    }
    catch {

    }
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
    <div className="flex flex-col items-center justify-center">
      <h2 className="text-xl mb-4 text-black">Verify Your Email</h2>
      <form onSubmit={handleSubmit} className="flex flex-col items-center space-y-4">
        <div className="flex items-center justify-center space-x-2">
          <InputOTP maxLength={6} pattern={REGEXP_ONLY_DIGITS_AND_CHARS} onChange={handleOtpChange} value={otp} onKeyDown={handleKeyDown} className="text-black">
            <InputOTPGroup className="text-black color-black">
              <InputOTPSlot index={0} classname="border-gray-500" />
              <InputOTPSlot index={1} className="border-gray-500" />
              <InputOTPSlot index={2} className="border-gray-500" />
            </InputOTPGroup>
            <InputOTPSeparator className="text-black" />
            <InputOTPGroup className="text-black">
              <InputOTPSlot index={3} className="border-gray-500" />
              <InputOTPSlot index={4} className="border-gray-500" />
              <InputOTPSlot index={5} className="border-gray-500" />
            </InputOTPGroup>
          </InputOTP>
        </div>
      </form >
    </div >
  );
};

export default VerifyEmailPage;
