import React from 'react';
import { Routes, Route } from 'react-router-dom';
import Login_Layout from './components/login/Login-Layout.jsx';
import Login_Email from './components/login/Login-Email.jsx';
import Login_Password from './components/login/Login-Password.jsx';
import Login_Update_Password from './components/login/Login-Update-Password.jsx';
// import Login_Verify from './components/login/Login-Verify.jsx'

import Form_Layout from './components/form/Form-Layout.jsx'
import Form from './components/form/Form.jsx'
// import Form from './Form.jsx'

import './App.css'

function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login_Layout />}>
        <Route path="/login" element={<Login_Email />} />
        <Route path="/login/password" element={<Login_Password />} />
        <Route path="/login/update-password" element={<Login_Update_Password />} />
      </Route>
      <Route path="/form" element={<Form />} />
    </Routes >
  )
}

export default App
