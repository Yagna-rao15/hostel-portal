import React from 'react';
import { Outlet } from 'react-router-dom';
import SvnitLogo from '../../assets/svnit-logo.png';

function Login_Layout() {


  return (
    <div>
      <header>
        <img src={SvnitLogo} className='logo react' />
        <h1 className="text-black"><strong>SVNIT SURAT</strong></h1>
        <p className="text-black"><strong>Hostel Section</strong></p>
        <br></br>
      </header>
      <main>
        <Outlet />
      </main>

    </div>
  );
}

export default Login_Layout;

