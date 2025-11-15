package hub

type MessageType string

const (
	DirectMessage MessageType = "direct_message"
	TeamBroadcast MessageType = "team_broadcast"
)

type Message[T any] struct {
	Type    MessageType `json:"type"`
	Payload T           `json:"payload"`
}

func NewMessage[T any](msgType MessageType, payload T) *Message[T] {
	return &Message[T]{
		Type:    msgType,
		Payload: payload,
	}
}
