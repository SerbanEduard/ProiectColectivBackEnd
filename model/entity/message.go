package entity

import (
	"errors"
	"strings"
	"time"
)

const (
	BadConversationKey = "bad conversation key"
)

type Message struct {
	ID              string    `json:"id"`
	SenderID        string    `json:"senderId"`
	SentAt          time.Time `json:"timestamp"`
	ConversationKey string    `json:"convKey,omitempty"`
	TeamID          string    `json:"teamId,omitempty"`
	TextContent     string    `json:"textContent"`
}

func NewMessage(id, senderId, convKey, teamId, textContent string) *Message {
	return &Message{
		ID:              id,
		SenderID:        senderId,
		SentAt:          time.Now().UTC(),
		ConversationKey: convKey,
		TeamID:          teamId,
		TextContent:     textContent,
	}
}

func GetConversationKey(user1Id, user2Id string) string {
	if user1Id < user2Id {
		return user1Id + "_" + user2Id
	}
	return user2Id + "_" + user1Id
}

func GetReceiverIdFromKey(senderId, conversationKey string) (string, error) {
	parts := strings.Split(conversationKey, "_")
	if len(parts) != 2 {
		return "", errors.New(BadConversationKey)
	}

	if parts[0] == senderId {
		return parts[1], nil
	}
	return parts[0], nil
}
