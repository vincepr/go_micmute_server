/*
*
*
 */

package main

import (
	"encoding/json"
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
	// create a ticker that triggers the ping signal to check if connection is alive
	ticker := time.NewTicker(pingInterval)

	// gracefully close and cleanup
	defer func() {
		ticker.Stop()
		c.manager.removeReceiverClient(c)
	}()

	for {
		select {
		case event, ok := <-c.eventQueue:
			// we check for End of Conneciton:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil{
					log.Println("connection CloseMessage sending failed with:", err)
				}
				return // we received the Close Signal -> we exit
			}
			// pack event we want to send:
			data, err := json.Marshal(event)
			if err != nil {
				log.Println("Failed to json.Marshal", err)
			}
			// send regular event to connection:
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err !=nil {
				log.Println("Failed Writing to Channel:", err)
			}
			log.Println("TODO: Remove: Send Event Successfully.")

		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("no pong received in time:", err)
				return	// got no answer -> we assume conection died
			}
		}
	}
}

// setup next Ping timer
func (c *ReceiverClient) pongReceiverHandler(pongMsg string) error {
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}