/*
*
*
 */

package main

import "encoding/json"

// Wrapper for different Event Types. All Incoming and Outgoing Messages use this as `encoding`
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Signature for Diferent Event-Types
//type EventHandler func(event Event, c *Client) error

// All supported Event Types:
const (
	SignalToReceiver  = "change_status"  	// Controller -> Server -> Receiver (ex. Volume_Up, Toggle_Mic)
	SignalToServer = "status_changed" 		// C# > Server (ex. Disconnecting)
)

type EventToReceiver struct {
	Id string `json:"id"`
	
}