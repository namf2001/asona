package websocket

import (
	"context"

	"github.com/redis/go-redis/v9"

	"asona/internal/model"
)

// Room represents a websocket room.
type Room struct {
	id         string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *model.WebsocketMessage
	redis      *redis.Client
}

// NewRoom creates a new Room.
func NewRoom(id string, rds *redis.Client) *Room {
	return &Room{
		id:         id,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *model.WebsocketMessage),
		redis:      rds,
	}
}

// Run runs the room, accepting various requests.
func (r *Room) Run(ctx context.Context) {
	go r.subscribeToRoomMessages(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case client := <-r.register:
			r.registerClient(client)
		case client := <-r.unregister:
			r.unregisterClient(client)
		case message := <-r.broadcast:
			r.publishMessage(ctx, message.Encode())
		}
	}
}

// ID returns the ID of the room.
func (r *Room) ID() string {
	return r.id
}

// registerClient adds the client to the room.
func (r *Room) registerClient(client *Client) {
	r.clients[client] = true
}

// unregisterClient removes the client from the room.
func (r *Room) unregisterClient(client *Client) {
	delete(r.clients, client)
}

// broadcastToClients sends the given message to all members in the room.
func (r *Room) broadcastToClients(message []byte) {
	for client := range r.clients {
		client.send <- message
	}
}

// publishMessage publishes the message to all clients subscribing to the room via Redis.
func (r *Room) publishMessage(ctx context.Context, message []byte) {
	if err := r.redis.Publish(ctx, r.id, message).Err(); err != nil {
		// Log error but don't block
	}
}

// subscribeToRoomMessages subscribes to messages in this room via Redis pub/sub.
func (r *Room) subscribeToRoomMessages(ctx context.Context) {
	pubsub := r.redis.Subscribe(ctx, r.id)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			r.broadcastToClients([]byte(msg.Payload))
		}
	}
}
