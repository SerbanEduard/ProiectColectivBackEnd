package dto

import "github.com/SerbanEduard/ProiectColectivBackEnd/model"

type TeamRequest struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	IsPublic    bool                  `json:"ispublic"`
	UserId      string                `json:"userid"`
	TeamTopic   model.TopicOfInterest `json:"teamtopic"`
}

type TeamResponse struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	IsPublic    bool                  `json:"ispublic"`
	UserId      string                `json:"userid"`
	TeamTopic   model.TopicOfInterest `json:"teamtopic"`
}

func NewTeamRequest(name, desc string, isPublic bool, userid string, teamtopic model.TopicOfInterest) *TeamRequest {
	return &TeamRequest{
		Name:        name,
		Description: desc,
		IsPublic:    isPublic,
		UserId:      userid,
		TeamTopic:   teamtopic,
	}
}

func NewTeamResponse(name, desc string, isPublic bool, userid string, teamtopic model.TopicOfInterest) *TeamResponse {
	return &TeamResponse{
		Name:        name,
		Description: desc,
		IsPublic:    isPublic,
		UserId:      userid,
		TeamTopic:   teamtopic,
	}
}

type UserToTeamRequest struct {
	UserID string `json:"userId" binding:"required"`
	TeamID string `json:"teamId" binding:"required"`
}

func NewUserToTeamRequest(userId, teamId string) *UserToTeamRequest {
	return &UserToTeamRequest{
		UserID: userId,
		TeamID: teamId,
	}
}
