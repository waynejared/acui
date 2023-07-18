import logo from './logo.svg';
import './App.css';
import Amplify, { API } from 'aws-amplify'
import React, { useEffect, useState } from 'react'

const myAPI = "RestLmabda"
const path = '/requestaccess'; 

const App = () => {
  const [input, setInput] = useState("")
  const [customers, setCustomers] = useState([])

  //Function to fetch from our backend and update customers array
  function setLockStatus(e, action) {
    let serialNumber = e.input
    let doorAction = action
    fetch("https://0uspubbub7.execute-api.us-west-2.amazonaws.com/staging/action/lockstatus/", {
      method: 'POST',
      body: JSON.stringify ({
          "serialnumber": serialNumber,
          "command": {
            "name": "lockstatus",
            "path": "/action/lockstatus",
            "value": doorAction
          }
      }),
      headers: {
        'Content-type': 'application/json; charset=UTF-8',
        'Access-Control-Request-Method':'POST',
      }
      })
      .then((response) => response.json())
      .then((data) => {
         console.log(data)
       })
       .catch(error => {
         console.log(error)
       })
  }

  return (
    
    <div className="App">
      <h1>Basic Door Control</h1>
      <div>
          <input placeholder="Serial Number" type="text" value={input} onChange={(e) => setInput(e.target.value)}/>      
      </div>
      <br/>
      <br/>
      <button onClick={() => setLockStatus({input},"unlocked")}>Unlock</button>&nbsp;&nbsp;
      <button onClick={() => setLockStatus({input},"locked")}>Lock</button>&nbsp;&nbsp;
      <button onClick={() => setLockStatus({input},"normal")}>Schedule</button>

      <h2 style={{visibility: customers.length > 0 ? 'visible' : 'hidden' }}>Response</h2>
      {
       customers.map((thisCustomer, index) => {
         return (
        <div key={thisCustomer.customerId}>
          <span><b>CustomerId:</b> {thisCustomer.customerId} - <b>CustomerName</b>: {thisCustomer.customerName}</span>
        </div>)
       })
      }
    </div>
  )
}

export default App;
/*import logo from "./logo.svg";
import "@aws-amplify/ui-react/styles.css";
import {
  withAuthenticator,
  Button,
  Heading,
  Image,
  View,
  Card,
} from "@aws-amplify/ui-react";

function App({ signOut }) {
  return (
    <View className="App">
      <Card>
        <Image src={logo} className="App-logo" alt="logo" />
        <Heading level={1}>We now have Auth!</Heading>
      </Card>
      <Button onClick={signOut}>Sign Out</Button>
    </View>
  );
}

export default withAuthenticator(App);
*/