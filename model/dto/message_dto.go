package dto

import "time"

type DirectMessageRequest struct {
	SenderID    string `json:"senderId"`
	ReceiverID  string `json:"receiverId"`
	TextContent string `json:"textContent"`
}

type TeamMessageRequest struct {
	SenderID    string `json:"senderId"`
	TeamId      string `json:"teamId"`
	TextContent string `json:"textContent"`
}

func NewDirectMessageRequest(senderId, receiverId, textContent string) *DirectMessageRequest {
	return &DirectMessageRequest{
		SenderID:    senderId,
		ReceiverID:  receiverId,
		TextContent: textContent,
	}
}

func NewTeamMessageRequest(senderId, teamId, textContent string) *TeamMessageRequest {
	return &TeamMessageRequest{
		SenderID:    senderId,
		TeamId:      teamId,
		TextContent: textContent,
	}
}

type MessageDTO struct {
	ID          string    `json:"id"`
	Sender      SenderDTO `json:"sender"`
	SentAt      string    `json:"sentAt"`
	ReceiverID  string    `json:"receiverId,omitempty"`
	TeamID      string    `json:"teamId,omitempty"`
	TextContent string    `json:"textContent"`
}

func NewMessageDTO(id, receiverId, teamId, textContent string, sentAt time.Time, sender SenderDTO) *MessageDTO {
	return &MessageDTO{
		ID:          id,
		Sender:      sender,
		SentAt:      sentAt.Format(time.RFC3339),
		ReceiverID:  receiverId,
		TeamID:      teamId,
		TextContent: textContent,
	}
}
