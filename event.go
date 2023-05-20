/*
*
*
 */

package main

import (
	"encoding/json"
)

// Wrapper for different Event Types. Corresponds to MEssages sent over the WebSocket.
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Signature to easily extend Different Event-Types later
type EventHandler func(event Event, c *ReceiverClient) error

// All supported Event Types:
const (
	EventSignalToReceiver  = "receiver_signal"  	// Controller -> Server -> Receiver (ex. Volume_Up, Toggle_Mic)
)

type SignalToReceiver struct {
	Id string `json:"id"`	
}