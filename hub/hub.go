package hub

import (
	"log"
	"sync"
)

type Hub[T any] struct {
	// The clients connected to this hub
	clients map[string]*Client[T]
	mu      sync.RWMutex
}

func NewHub[T any]() *Hub[T] {
	return &Hub[T]{
		clients: make(map[string]*Client[T]),
	}
}

func (h *Hub[T]) Register(client *Client[T]) {
	h.mu.Lock()
	h.clients[client.ClientID] = client
	h.mu.Unlock()

	// Start the read and write pump
	go h.readPump(client)
	go h.writePump(client)
}

func (h *Hub[T]) Unregister(client *Client[T]) {
	h.mu.Lock()
	_, ok := h.clients[client.ClientID]
	if ok {
		delete(h.clients, client.ClientID)
	}
	h.mu.Unlock()

	if ok {
		close(client.outbound)
		err := client.Conn.Close()
		if err != nil {
			log.Printf("Failed to close the connection for client with ID = %s", client.ClientID)
			return
		}
	}
}

func (h *Hub[T]) Send(clientID string, msg T) {
	h.mu.Lock()
	client, ok := h.clients[clientID]
	h.mu.Unlock()

	if !ok {
		// Client is offline
		return
	}

	select {
	case client.outbound <- msg:
		// Sent to outbound channel
	default:
		// Channel is full
	}
}

func (h *Hub[T]) SendMany(clientIDs []string, msg T) {
	for _, id := range clientIDs {
		h.Send(id, msg)
	}
}

// writePump continuously reads from the outbound channel and writes to the WebSocket
func (h *Hub[T]) writePump(client *Client[T]) {
	defer func() {
		// Unregister the client on exit
		h.Unregister(client)
	}()

	for msg := range client.outbound {
		if err := client.Conn.WriteJSON(msg); err != nil {
			// Client disconnected
			return
		}
	}
}

// readPump continuously checks for disconnection
func (h *Hub[T]) readPump(client *Client[T]) {
	defer func() {
		h.Unregister(client)
	}()

	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			// WebSocket sent a disconnect message
			return
		}
	}
}
