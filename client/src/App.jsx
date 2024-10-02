import React from 'react';
import { Routes, Route } from 'react-router-dom';
import Login_Layout from './components/login/Login-Layout.jsx';
import Login_Email from './components/login/Login-Email.jsx';
// import Login_Password from './components/login/Login-Password.jsx';
// import Login_Verify from './components/login/Login-Verify.jsx'
import './App.css'

function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login_Layout />}>
        <Route path="/login" element={<Login_Email />} />
      </Route>
    </Routes>
  )
}

export default App
