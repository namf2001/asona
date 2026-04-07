package model

import (
	"encoding/json"
	"log"
)

// ReceivedMessage represents a received websocket message from client.
type ReceivedMessage struct {
	Action  string `json:"action"`
	Room    string `json:"room"`
	Message any    `json:"message,omitempty"`
}

// WebsocketMessage represents an emitted websocket message to client.
type WebsocketMessage struct {
	Action string `json:"action"`
	Data   any    `json:"data,omitempty"`
}

// Encode turns the message into a byte array.
func (m *WebsocketMessage) Encode() []byte {
	encoding, err := json.Marshal(m)
	if err != nil {
		log.Printf("websocket: failed to encode message: %v", err)
		return nil
	}
	return encoding
}
