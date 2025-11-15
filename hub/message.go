package hub

type MessageType string

const (
	DirectMessage MessageType = "direct_message"
	TeamBroadcast MessageType = "team_broadcast"
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

func NewMessage(msgType MessageType, payload interface{}) *Message {
	return &Message{
		Type:    msgType,
		Payload: payload,
	}
}
