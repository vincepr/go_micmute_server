/*
*
*
 */

package main

import (
	"encoding/json"
	"fmt"
)

// Wrapper for different Event Types. Corresponds to MEssages sent over the WebSocket.
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Signature to easily extend Different Event-Types later
type EventHandler func(event Event, c *ReceiverClient) error

// string-Identifiers for all supported Event Types:
const (
	eventSignalToReceiver  = "receiver_signal"  	// Controller -> Server -> Receiver (ex. Volume_Up, Toggle_Mic)
)

// string-Identifiers for all supported Signals (that way we can call out wrong use of the api)
var supportedSignalsMap = map[string]bool{
	"vol_down": true,
	"vol_up": true,
	"vol_toggle": true,
	"mic_down": true,
	"mic_up": true,
	"mic_toggle": true,
}

type SignalToReceiver struct {
	Signal string `json:"signal"`	
}

func NewSignalToReceiver(signal_str string) (*Event, error) {
	if _, ok := supportedSignalsMap[signal_str]; !ok{
		return nil, fmt.Errorf("unsupported Signal-identifier: %v", signal_str)
	}
	// build outgoing Signal:
	var sig = SignalToReceiver{ Signal: signal_str }
	data, err := json.Marshal(sig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal NeWSignalToReceiver: %v", err)
	}
	// wrap it in the Event
	return &Event{
		Type: eventSignalToReceiver,
		Payload: data,
	}, nil
}