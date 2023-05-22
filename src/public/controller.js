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
        "signal": eventname
    }
    fetch("controller", {
        method: "post",
        body: JSON.stringify(formData),
        mode: "cors",
    }).then((response) => {
        if(response.ok) console.log(`${eventname} sent successfully`);
        else throw 'unauthorized, use the username and password you get when running the MicMute.exe';
    }).catch((err) => {customAlert(err)});
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

/*
*   Cusom Alert to display error messages and remove them afer a timeout
*/

function customAlert(textToDisplay) {
    let x = document.getElementById("customAlert");
    document.getElementById("customAlertText").innerHTML = textToDisplay;
    x.className = "show";
    setTimeout(function(){ x.className = x.className.replace("show", ""); }, 3000);
}