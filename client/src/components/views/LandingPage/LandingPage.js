import React from 'react'
import { FaCoffee, FaReact } from "react-icons/fa";

function LandingPage() {
    return (
        <>
        <div className="app">
    
        <FaCoffee style={{ fontSize: '6.5rem' }}/>

            <span style={{ fontSize: '2rem' }}>Let's Start Coding!</span>
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
