package hub

import "github.com/gorilla/websocket"

type Client[T any] struct {
	ClientID string
	Conn     *websocket.Conn

	// Channel for sending messages to the client
	outbound chan T
}

func NewClient[T any](clientID string, conn *websocket.Conn) *Client[T] {
	return &Client[T]{
		ClientID: clientID,
		Conn:     conn,
		outbound: make(chan T),
	}
}
