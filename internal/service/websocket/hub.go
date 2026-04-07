package websocket

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"

	"asona/internal/model"
)

// Hub contains all rooms and clients.
type Hub struct {
	ctx         context.Context
	clients     map[*Client]bool
	register    chan *Client
	unregister  chan *Client
	broadcast   chan []byte
	rooms       map[*Room]bool
	redisClient *redis.Client
	mu          sync.RWMutex
}

// NewHub creates a new Hub.
func NewHub(ctx context.Context, redisClient *redis.Client) *Hub {
	return &Hub{
		ctx:         ctx,
		clients:     make(map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan []byte),
		rooms:       make(map[*Room]bool),
		redisClient: redisClient,
	}
}

// Run runs the hub, accepting various requests.
func (h *Hub) Run() {
	for {
		select {
		case <-h.ctx.Done():
			return
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case message := <-h.broadcast:
			h.broadcastToClients(message)
		}
	}
}

// Register registers a client to the hub.
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// registerClient adds the client to the hub and starts pumps.
func (h *Hub) registerClient(client *Client) {
	go client.WritePump(h.ctx)
	go client.ReadPump(h.ctx)

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()
}

// unregisterClient removes the client from the hub.
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	delete(h.clients, client)
	h.mu.Unlock()
}

// broadcastToClients sends the given message to all connected clients.
func (h *Hub) broadcastToClients(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			// Client buffer full, skip
		}
	}
}

// BroadcastToRoom sends the given message to all clients in the specified room.
func (h *Hub) BroadcastToRoom(message *model.WebsocketMessage, roomID string) {
	h.mu.RLock()
	room := h.findRoomByIDLocked(roomID)
	h.mu.RUnlock()

	if room != nil {
		room.publishMessage(h.ctx, message.Encode())
	}
}

// findRoomByID finds a room by its ID.
func (h *Hub) findRoomByID(id string) *Room {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.findRoomByIDLocked(id)
}

// findRoomByIDLocked finds a room by its ID (must hold lock).
func (h *Hub) findRoomByIDLocked(id string) *Room {
	for room := range h.rooms {
		if room.ID() == id {
			return room
		}
	}
	return nil
}

// createRoom creates a new room and starts it.
func (h *Hub) createRoom(id string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Double check if room was created while waiting for lock
	if room := h.findRoomByIDLocked(id); room != nil {
		return room
	}

	room := NewRoom(id, h.redisClient)
	go room.Run(h.ctx)
	h.rooms[room] = true

	return room
}
