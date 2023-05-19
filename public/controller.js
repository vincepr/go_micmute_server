/*
*   Setting globals up and connecting onclick-/form-events
*/

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
    let formData = {
        "username": document.getElementById("username").value,
        "password": document.getElementById("password").value,
        "event": eventname
    }
    fetch("controller", {
        method: "post",
        body: JSON.stringify(formData),
        mode: "cors",
    }).then((response) => {
        if(response.ok) return response.json();
        else throw 'unauthorized';
    }).then((data) => {
        console.log(data)
        //
    }).catch((err) => {alert(err)});
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