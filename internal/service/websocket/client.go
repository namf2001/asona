package websocket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"

	"asona/internal/constants"
	"asona/internal/model"
)

var newline = []byte{'\n'}

// Client represents a websocket client at the server.
type Client struct {
	id    string
	conn  *websocket.Conn
	hub   *Hub
	send  chan []byte
	rooms map[*Room]bool
}

// NewClient creates a new Client.
func NewClient(id string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		id:    id,
		conn:  conn,
		hub:   hub,
		send:  make(chan []byte, 256),
		rooms: make(map[*Room]bool),
	}
}

// ID returns the client ID.
func (c *Client) ID() string {
	return c.id
}

// ReadPump pumps messages from the websocket connection to the hub.
func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		c.disconnect()
	}()

	c.conn.SetReadLimit(constants.MaxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(constants.PongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(constants.PongWait))
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, jsonMessage, err := c.conn.ReadMessage()
			if err != nil {
				return
			}
			c.handleNewMessage(jsonMessage)
		}
	}
}

// WritePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump(ctx context.Context) {
	ticker := time.NewTicker(constants.PingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(constants.WriteWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Attach queued messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-c.send)
			}

			if err = w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(constants.WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// disconnect disconnects the client from all rooms and the hub.
func (c *Client) disconnect() {
	c.hub.unregister <- c
	for room := range c.rooms {
		room.unregister <- c
	}
	close(c.send)
	_ = c.conn.Close()
}

// handleNewMessage handles incoming messages from the client.
func (c *Client) handleNewMessage(jsonMessage []byte) {
	var message model.ReceivedMessage
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		c.sendError("invalid message format")
		return
	}

	switch message.Action {
	case constants.JoinRoomAction:
		c.handleJoinRoom(message)
	case constants.LeaveRoomAction:
		c.handleLeaveRoom(message)
	case constants.StartTyping:
		c.handleTypingEvent(message, constants.AddToTypingAction)
	case constants.StopTyping:
		c.handleTypingEvent(message, constants.RemoveFromTypingAction)
	case constants.ToggleOnline:
		c.handleToggleOnline(true)
	case constants.ToggleOffline:
		c.handleToggleOnline(false)
	default:
		c.sendError("unknown action")
	}
}

// handleJoinRoom joins the given room.
func (c *Client) handleJoinRoom(message model.ReceivedMessage) {
	roomID := message.Room
	if roomID == "" {
		c.sendError("room id is required")
		return
	}

	room := c.hub.findRoomByID(roomID)
	if room == nil {
		room = c.hub.createRoom(roomID)
	}

	c.rooms[room] = true
	room.register <- c
}

// handleLeaveRoom leaves the given room.
func (c *Client) handleLeaveRoom(message model.ReceivedMessage) {
	room := c.hub.findRoomByID(message.Room)
	if room != nil {
		delete(c.rooms, room)
		room.unregister <- c
	}
}

// handleTypingEvent emits typing status to the room.
func (c *Client) handleTypingEvent(message model.ReceivedMessage, action string) {
	roomID := message.Room
	if room := c.hub.findRoomByID(roomID); room != nil {
		msg := model.WebsocketMessage{
			Action: action,
			Data:   map[string]any{"userId": c.id, "message": message.Message},
		}
		room.broadcast <- &msg
	}
}

// handleToggleOnline handles online/offline status toggle.
func (c *Client) handleToggleOnline(isOnline bool) {
	action := constants.ToggleOfflineEmission
	if isOnline {
		action = constants.ToggleOnlineEmission
	}

	// Broadcast to user's room
	if room := c.hub.findRoomByID(c.id); room != nil {
		msg := model.WebsocketMessage{
			Action: action,
			Data:   map[string]any{"userId": c.id},
		}
		room.broadcast <- &msg
	}
}

// sendError sends an error message to the client.
func (c *Client) sendError(errMsg string) {
	msg := model.WebsocketMessage{
		Action: constants.ErrorAction,
		Data:   map[string]any{"error": errMsg},
	}
	c.send <- msg.Encode()
}
