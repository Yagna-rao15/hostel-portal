// import React from "react";
import { Input } from "@/components/ui/input"

function Login() {
  return (
    <div className="head">
      <image src="./assets/svnot-logo.png" alt="SVNIT-LOGO" />
      <h1 className="text-black"><b>SVNIT SURAT</b></h1>
      <p className="text-black"><b>Hostel Section</b></p>
      <br></br>
      <Input type="email" placeholder="Type your email/username" />
      <br></br>
      <Input type="password" placeholder="Password" />
    </div>
  )
}

export default Login
