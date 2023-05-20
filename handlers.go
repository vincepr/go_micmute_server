/*
*		Handlers for the different API Routes like "/login" "/ws"
*
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
*		Controllers are Clients(Website) who want to controll the others Microphone and Volume settings
*		They send controllerRequests via HTML Post-Requests, including validation and the signal type.
*/

// Controller tries to send these signals to turn Volume Up, Mute the mic etc...
type controllerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Signal string `json:"signal"`
}

func (m *Manager)ControllerRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req controllerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check if username and password match up
	if !m.isValidUsernamePw(req.Username, req.Password) {
		http.Error(w, "failed authorisation", http.StatusUnauthorized)
		return
	}
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
	if m.isUsernameInUse(otp.Username) {
		http.Error(w, "username already in use", http.StatusNotAcceptable)
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
	m.addReceiverClient(client)
	go client.getEvents()
	go client.sendEvents()
}
