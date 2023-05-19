/*
*
*
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// pongWait is how long we wait, before assuming the client is dead and cleanup
	// pingInterval is when we send the next ping, that our client has to answer
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
	// in bytes. Maximum size of one Message.
	msgMaxSize int64 = 512
)


/*
*		Client that acts as the Controller (webbrowser) -> sending the Mute/Volume-up signals
*/

type ControllerClient struct {
	conn           *websocket.Conn
	manager        *Manager   // ref to manager responsible for this connection
	eventQueue     chan Event // all yet to handle Events for this WebSocket (channel blocks so async save)
	targetUsername string
	targetPassword string
}

func NewControllerClient(conn *websocket.Conn, manager *Manager, otp OTP) *ControllerClient {
	return &ControllerClient{
		conn:           conn,
		manager:        manager,
		eventQueue:     make(chan Event),
		targetUsername: otp.Username,
		targetPassword: otp.Password,
	}
}

// from Browser -> Server
func (c *ControllerClient) getEvents() {
	defer func() {
		c.manager.removeControllerClient(c)
	}()
	// Setup auto Disconect(using ping/pong) and maxMessageSize
	c.conn.SetReadLimit(msgMaxSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println("could not set timeout:", err)
		return
	}
	c.conn.SetPongHandler(c.pongControllerHandler)

	// Handling incoming Events
	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Connection Has Closed Unexpected: ", err)
			}
			break
		}
		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("Error json.Unmarshal: %v", err)
			break // TODO: maybe handle this gracefully? request sending again etc...
		}
		// Route the Event
		// if err := c.manager.routeEvent(request, c); err != nil {
		// 	log.Println("Error routeEvent:", err)
		// }
	}
}

// from this Server -> Webbrowser
func (c *ControllerClient) sendEvents() {
	fmt.Println("controller sending")
}

// setup next Ping timer
func (c *ControllerClient) pongControllerHandler(pongMsg string) error {
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}

/*
* 		Client that acts as the Receiver (C# APP) -> receiving and acting on the Mute/Volume-up signals
*/

type ReceiverClient struct {
	conn       *websocket.Conn
	manager    *Manager   // ref to manager responsible for this connection
	eventQueue chan Event // all yet to handle Events for this WebSocket (channel blocks so async save)
	username   string
	password   string
}

func NewReceiverClient(conn *websocket.Conn, manager *Manager, otp OTP) *ReceiverClient {
	return &ReceiverClient{
		conn:       conn,
		manager:    manager,
		eventQueue: make(chan Event),
		username:   otp.Username,
		password:   otp.Password,
	}
}

// from c#-application -> this Server
func (c *ReceiverClient) getEvents() {
	defer func() {
		c.manager.removeReceiverClient(c)
	}()
	// Setup auto Disconect(using ping/pong) and maxMessageSize:
	c.conn.SetReadLimit(msgMaxSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println("could not set timeout:", err)
		return
	}
	c.conn.SetPongHandler(c.pongReceiverHandler)

	// Handling incoming Events:
	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Connection Has Closed Unexpected: ", err)
			}
			break
		}
		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("Error json.Unmarshal: %v", err)
			break // TODO: maybe handle this gracefully? request sending again etc...
		}
		// Route the Event
		// if err := c.manager.routeEvent(request, c); err != nil {
		// 	log.Println("Error routeEvent:", err)
		// }
	}
}

// from this Server -> c#-application
func (c *ReceiverClient) sendEvents() {
	fmt.Println("controller sending")
}

// setup next Ping timer
func (c *ReceiverClient) pongReceiverHandler(pongMsg string) error {
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}