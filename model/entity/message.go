package entity

import "time"

type Message struct {
	ID              string    `json:"id"`
	SenderID        string    `json:"senderId"`
	SentAt          time.Time `json:"timestamp"`
	ConversationKey string    `json:"convKey,omitempty"`
	TeamId          string    `json:"teamId,omitempty"`
	TextContent     string    `json:"textContent"`
}

func NewMessage(id, senderId, convKey, teamId, textContent string) *Message {
	return &Message{
		ID:              id,
		SenderID:        senderId,
		SentAt:          time.Now().UTC(),
		ConversationKey: convKey,
		TeamId:          teamId,
		TextContent:     textContent,
	}
}

func GetConversationKey(user1Id, user2Id string) string {
	if user1Id < user2Id {
		return user1Id + "_" + user2Id
	}
	return user2Id + "_" + user1Id
}
