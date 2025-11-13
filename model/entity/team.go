package entity

import "github.com/SerbanEduard/ProiectColectivBackEnd/model"

type Team struct {
	Id          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	IsPublic    bool                  `json:"ispublic"`
	UsersIds    []string              `json:"users"`
	TeamTopic   model.TopicOfInterest `json:"teamtopic"`
}

func NewTeam(id, name, desc string, isPublic bool, Users []string, topic model.TopicOfInterest) *Team {
	return &Team{
		Id:          id,
		Name:        name,
		Description: desc,
		IsPublic:    isPublic,
		UsersIds: func() []string {
			if Users == nil {
				return []string{}
			}
			return Users
		}(),
		TeamTopic: topic,
	}
}
