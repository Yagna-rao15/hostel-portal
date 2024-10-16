import React from 'react';
import SvnitLogo from '../../assets/svnit-logo.png';

function Login_Layout() {


  return (
    <div>
      <div className="flex items-center justify-between p-0 gap-4">
        <div className="flex-shrink-0">
          <img src={SvnitLogo} className="h-16 w-auto" alt="SVNIT Logo" />
        </div>
        <div className="text-center flex-grow">
          <h2 className="text-black text-4xl font-bold ">Sardar Vallabhbhai</h2>
          <h3 className="text-black text-xl ">National Institute of Technology</h3>
        </div>
      </div>
    </div>
  );
}

export default Login_Layout;


