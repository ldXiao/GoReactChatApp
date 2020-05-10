import React from 'react'
import {Icon} from 'antd';

function Footer() {
    return (
        <div style={{
            height: '50px', display: 'flex',
            flexDirection: 'column', alignItems: 'center',
            justifyContent: 'center', fontSize:'1rem'
        }}>
           <p> Enjoy  <Icon type="smile" /></p>
        </div>
    )
}

export default Footer
