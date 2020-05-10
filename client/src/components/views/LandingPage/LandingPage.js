import React from 'react'
import { FaCoffee, FaReact } from "react-icons/fa";
import {GO} from "../Icons/Icons";

import "../../../devicon.css";

function LandingPage() {
    return (
        <>
        <div className="app">
        <div>
        <FaReact style={{color:"#49c" ,fontSize: '6.5rem' }}/>
        <GO width={100} fill="#49c" />
        </div>
       

            <span style={{ fontSize: '2rem' }}>Let's Start Chatting!</span>
        </div>
        <div style={{
            height: '20px', display: 'flex',
            flexDirection: 'column', alignItems: 'center',
            justifyContent: 'center', fontSize:'1rem'}}>
            Frontend addapted from Boiler Plate by John Ahn</div>
        <br/>
        <div 
        style={{
            height: '20px', display: 'flex',
            flexDirection: 'column', alignItems: 'center',
            justifyContent: 'center', fontSize:'1rem'}}>
            Backend Powered by Golang wirtten by Lind Xiao</div>
        </>
    )
}

export default LandingPage
