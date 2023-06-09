/*		Handlers for the different API Routes like "/login" "/ws"
*		- Controllers use stateless HTTP at the moment
*		- Receivers use HTTP for login then stay connected on a WebSocket
*/
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrades the incoming HTTP(S) request to a Websocket
var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

/*
*		Controllers are Clients (from Website) who want to controll the others Microphone and Volume settings
*		They send controllerRequests via HTML Post-Requests, including validation and the signal type.
 */

// Controller tries to send these signals to turn Volume Up, Mute the mic etc...
type controllerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Signal   string `json:"signal"`
}

func (m *Manager) ControllerRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req controllerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check if username and password match up
	receiver, ok := m.isValidUsernamePw(req.Username, req.Password);
	if !ok {
		http.Error(w, "failed authorisation", http.StatusUnauthorized)
		return
	}
	// build the event we want to send:
	event, err := NewSignalToReceiver(req.Signal)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// push the event into the receivers queue:
	receiver.eventQueue <- *event

	w.WriteHeader(http.StatusOK)
}

/*
*		Receivers are Clients who want their Microphone/Volume controlled by Signals sent to them
*		- They Send credentials via http, get a OTP back, use OTP to try to Update to WS
*		- And then stay connected via WS, continously listening for Events/Signals
*		- They must respond to continous pings with pongs otherwise connection drops
 */

type receiverLoginRequest struct {
	Apikey   string `json:"apikey"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type responseReceiver struct {
	OTP string `json:"otp"`
}

func (m *Manager) loginReceiverHandler(w http.ResponseWriter, r *http.Request) {
	var req receiverLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Apikey == "noauth" {
		otp := m.otps.NewOTP(req.Username, req.Password)
		resp := responseReceiver{
			OTP: otp.Key,
		}
		data, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}
	http.Error(w, "wrong api key", http.StatusUnauthorized)
}

// Checks OTP from HTTP request, then upgrades to WebSocket-connection and serves that.
func (m *Manager) serveReceiverHandler(w http.ResponseWriter, r *http.Request) {
	// check for valid OTP key:
	otp_key := r.URL.Query().Get("otp")
	if otp_key == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	otp, isOk := m.otps.VertifyOTP(otp_key)
	if !isOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Check if username already exists:
	if m.isUsernameTaken(otp.Username) {
		http.Error(w, "username already in use", http.StatusConflict)
		return
	}
	// upgrade the HTTP request to a Websocket Connection
	log.Println("New Receiver-connection upgrade request")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed upgrade to Websocket", err)
		return
	}
	client := NewReceiverClient(conn, m, otp)
	if ok := m.addReceiverClient(client); !ok {
		// 2 users could theoretically try to upgrade with same name so we close 2nd.
		client.conn.WriteMessage(websocket.CloseMessage, nil)
		return
	}
	go client.getEvents()
	go client.sendEvents()
}
