package constants

import "time"

// Subscribed Messages (from client)
const (
	JoinRoomAction  = "joinRoom"
	LeaveRoomAction = "leaveRoom"
	StartTyping     = "startTyping"
	StopTyping      = "stopTyping"
	ToggleOnline    = "toggleOnline"
	ToggleOffline   = "toggleOffline"
)

// Emitted Messages (to client)
const (
	NewMessageAction       = "new_message"
	AddToTypingAction      = "addToTyping"
	RemoveFromTypingAction = "removeFromTyping"
	ToggleOnlineEmission   = "toggle_online"
	ToggleOfflineEmission  = "toggle_offline"
	ErrorAction            = "error"
)

// WebSocket timeouts
const (
	// WriteWait is the max wait time when writing message to peer.
	WriteWait = 10 * time.Second

	// PongWait is the time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second

	// PingPeriod sends pings to peer with this period. Must be less than PongWait.
	PingPeriod = (PongWait * 9) / 10

	// MaxMessageSize is the maximum message size allowed from peer.
	MaxMessageSize = 10000
)
