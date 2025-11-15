package dto

type DirectMessagesRequest struct {
	SenderId    string `json:"senderId"`
	ReceiverId  string `json:"receiverId"`
	TextContent string `json:"textContent"`
}

type TeamMessagesRequest struct {
	SenderId    string `json:"senderId"`
	TeamId      string `json:"teamId"`
	TextContent string `json:"textContent"`
}

func NewDirectMessageRequest(senderId, receiverId, textContent string) *DirectMessagesRequest {
	return &DirectMessagesRequest{
		SenderId:    senderId,
		ReceiverId:  receiverId,
		TextContent: textContent,
	}
}

func NewTeamMessageRequest(senderId, teamId, textContent string) *TeamMessagesRequest {
	return &TeamMessagesRequest{
		SenderId:    senderId,
		TeamId:      teamId,
		TextContent: textContent,
	}
}
