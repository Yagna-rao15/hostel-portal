import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import SvnitLogo from '../../assets/svnit-logo.png';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const Login_Email = () => {
  const location = useLocation();
  const email = location.state?.email || '';
  const [name, setName] = useState('');
  const [room, setRoom] = useState('');
  const [mobile, setMobile] = useState('');
  const [description, setDescription] = useState('');
  const [type, setComplainType] = useState('');
  const [hostel, setHostel] = useState('');
  const [file, setFiles] = useState(null);
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);
  // const navigate = useNavigate();



  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const formData = new FormData();
      formData.append('email', email);
      formData.append('complainType', type);
      formData.append('hostel', hostel);
      formData.append('name', name);
      formData.append('room', room);
      formData.append('mobile', mobile);
      formData.append('description', description);
      formData.append('file', file); // Append the file

      const response = await fetch('http://localhost:8080/form', {
        method: 'POST',
        body: formData, // Send FormData instead of JSON
      });

      // Check if response is ok (status in the range 200-299)
      if (!response.ok) {
        setMessage('Failed to connect to server. Please try again later.');
        return;
      }

      const data = await response.json();

      // Validate response structure
      if (!data || typeof data !== 'object') {
        setMessage('Unexpected server response. Please try again.');
        return;
      } else {
        setMessage(JSON.stringify(data));
        // setMessage("Done submitting");
        // console.log(data);
      }

    } catch (error) {
      console.error('Error checking email:', error);
      setMessage('Server error. Please try again later.');
    } finally {
      setLoading(false); // Ensure loading is set to false in finally block
    }
  };

  return (
    <div>
      <div>
        <header>
          <div className="flex items-center justify-between p-0 gap-4">
            <div className="flex-shrink-0">
              <img src={SvnitLogo} className="h-16 w-auto" alt="SVNIT Logo" />
            </div>
            <div className="text-center flex-grow">
              <h2 className="text-black text-4xl font-bold ">Sardar Vallabhbhai</h2>
              <h3 className="text-black text-xl ">National Institute of Technology</h3>
            </div>
          </div>
        </header>
      </div>
      <br />
      <hr className="border-black" />
      <br />
      <form onSubmit={handleSubmit}>
        <div className='text-black font-bold text-left text-xl '>Complain:</div>
        <div className='text-gray-500 text-left text-base '>Register your Complain here</div>
        <br />
        <p className='text-black text-left text-base pb-2'>Complain Type</p>
        <Select required onValueChange={setComplainType}>
          <SelectTrigger className="w-96 text-black">
            <SelectValue placeholder="Type" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="General">General</SelectItem>
            <SelectItem value="Mess">Mess</SelectItem>
            <SelectItem value="Electricity">Electricity</SelectItem>
            <SelectItem value="Bathroom">Bathroom</SelectItem>
          </SelectContent>
        </Select>
        <br />
        <p className='text-black text-left text-base pb-2'>Select Hostel</p>
        <Select required onValueChange={setHostel}>
          <SelectTrigger className="w-96 text-black">
            <SelectValue placeholder="Select Hostel" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="Gajjar Bhavan">Gajjar Bhavan</SelectItem>
            <SelectItem value="Atal Bihari Vajpaye Bhavan">Atal Bihari Vajpaye Bhavan</SelectItem>
            <SelectItem value="Bhabha Bhavan">Bhabha Bhavan</SelectItem>
            <SelectItem value="Tagore Bhavan">Tagore Bhavan</SelectItem>
            <SelectItem value="Swamy Bhavan">Swamy Bhavan</SelectItem>
            <SelectItem value="Mother Teresa Bhavan">Mother Teresa Bhavan</SelectItem>
            <SelectItem value="Raman Bhavan">Raman Bhavan</SelectItem>
            <SelectItem value="Sarabhai Bhavan">Sarabhai Bhavan</SelectItem>
          </SelectContent>
        </Select>
        <br />
        <p className='text-black text-left text-base pb-2'>Enter Your Name</p>
        <Input
          className="text-black"
          type="text"
          placeholder="Enter your name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
        <br />
        <p className='text-black text-left text-base pb-2'>Enter Your Room Number and Floor</p>
        <Input
          className="text-black"
          type="text"
          placeholder="Type your room and floor"
          value={room}
          onChange={(e) => setRoom(e.target.value)}
          required
        />
        <br />
        <p className='text-black text-left text-base pb-2'>Enter Your Mobile Number</p>
        <Input
          className="text-black"
          type="text"
          placeholder="Enter your mobile number"
          value={mobile}
          onChange={(e) => setMobile(e.target.value)}
          required
        />
        <br />
        <p className='text-black text-left text-base pb-2'>Description</p>
        <Textarea
          placeholder="Type your description here."
          className='text-black'
          id="message-2"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />
        <br />
        <p className='text-black text-left text-base pb-2'>Related Pictures or Videos</p>
        <Input
          className="text-black text-left"
          type="file"
          onChange={(e) => setFiles(e.target.files[0])}
        />
        <br />
        <Button
          className="w-96 h-1em text-base"
          type="submit"
          disabled={loading}
        >
          {loading ? 'Verifying...' : 'Submit'}
        </Button>
        {message && <p className="text-red-500">{message}</p>}
      </form>
    </div>
  );
}

export default Login_Email;

