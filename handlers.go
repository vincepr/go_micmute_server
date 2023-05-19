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
 */

// TODO: we could Block/Whitelist all CrossSite requests not coming from "my" website when everything is running

type controllerLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type responseController struct {
	OTP string `json:"otp"`
}

func (m *Manager) loginControllerHandler(w http.ResponseWriter, r *http.Request) {
	var req controllerLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	otp := m.otps.NewOTP(req.Username, req.Password)
	resp := responseController{
		OTP: otp.Key,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (m *Manager) serveControllerHandler(w http.ResponseWriter, r *http.Request) {
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
	// upgrade the HTTP request to a Websocket Connection
	log.Println("New Controller-connection upgrade request")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed upgrade to Websocket", err)
		return
	}
	client := NewControllerClient(conn, m, otp)
	m.addControllerClient(client)
	go client.getEvents()
	go client.sendEvents()
}

/*
*		Receivers are Clients who want their Microphone/Volume controlled by Signals sent to them
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
