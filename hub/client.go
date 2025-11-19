package hub

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	clientOutboundBufferSize = 16
	webSocketReadBufferSize  = 1024
	webSocketWriteBufferSize = 1024
)

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
		outbound: make(chan T, clientOutboundBufferSize),
	}
}

func AcceptConnection(c *gin.Context) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  webSocketReadBufferSize,
		WriteBufferSize: webSocketWriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	return conn, err
}
