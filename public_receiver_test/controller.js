/*
*   Setting globals up and connecting onclick-/form-events
*/

var wsConn;

document.getElementById("login").addEventListener("submit", handleLogin);

document.getElementById("volDown")  .onclick = () => sendEvent("vol_down");
document.getElementById("volUp")    .onclick = () => sendEvent("vol_up");
document.getElementById("volToggle").onclick = () => sendEvent("vol_toggle");
document.getElementById("micDown")  .onclick = () => sendEvent("mic_down");
document.getElementById("micUp")    .onclick = () => sendEvent("mic_up");
document.getElementById("micToggle").onclick = () => sendEvent("mic_toggle");

/*
*   Handler functions:
*/

function sendEvent(eventname) {
    console.log(eventname);
}

function handleLogin(ev) {
    ev.preventDefault();
    let formData = {
        "username": document.getElementById("username").value,
        "password": document.getElementById("password").value,
    }
    // send request to the /login api endpoint
    fetch("login_controller", {
        method: "post",
        body: JSON.stringify(formData),
        mode: "cors",
    }).then((response) => {
        if(response.ok) return response.json();
        else throw 'unauthorized';
    }).then((data) => {
        console.log(data)
        console.log(data.otp)
        connectWebsocket(data.otp);
    }).catch((err) => {alert(err)});
}

/*
*   Functions dealing with the Websocket connection
*/

function connectWebsocket(oneTimePassword) {
    if (!window["WebSocket"]){
        alert("Browser not supporting websockets, change your browser.");
        return;
    }
    wsConn = new WebSocket("ws://"+ document.location.host + "/controller?otp="+oneTimePassword);
    setupWsEventHandlers();
}

function setupWsEventHandlers() {
    wsConn.onopen = () => {
        document.getElementById("connection-header").innerHTML = "Logged in - active Websocket connection.";
    }
    wsConn.onclose = () => {
        document.getElementById("connection-header").innerHTML = "Not Logged in - Websocket connection closed.";
    }
    wsConn.onerror = (ev) => {
        console.log("error with the websocket: "+ ev);
    }
    wsConn.onmessage = (ev) => receivedEvent(ev);
}

/*
*   Functions dealing with receiving Eventsignals over the Websocket connection:
*/

function receivedEvent(ev) {
    let eventData = JSON.parse(ev.data);
    let event = Object.assign(new Event, eventData);
    routeEvent(event);
}

function routeEvent(evt) {
    console.log(evt);
}

/*
*   Classes describing the JSON data coming from the Server
*/

// Wrapper other Event Types get wrapped into. (into the payload) 
class Event {
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}