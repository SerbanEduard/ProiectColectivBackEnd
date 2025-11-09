package dto

import (
	"time"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model"
)

type UpdateStatisticsRequest struct {
	TimeSpentOnApp string `json:"timeSpentOnApp"`
	TeamId         string `json:"teamId"`
	TimeSpentOnTeam string `json:"timeSpentOnTeam"`
}

type UpdateStatisticsResponse struct {
	UserId              string                  `json:"userId"`
	TotalTimeSpentOnApp time.Duration           `json:"totalTimeSpentOnApp"`
	TimeSpentOnTeams    []model.TimeSpentOnTeam `json:"timeSpentOnTeams"`
}

func NewUpdateStatisticsResponse(userId string, statistics *model.Statistics) *UpdateStatisticsResponse {
	return &UpdateStatisticsResponse{
		UserId:              userId,
		TotalTimeSpentOnApp: statistics.TotalTimeSpentOnApp,
		TimeSpentOnTeams:    statistics.TimeSpentOnTeams,
	}
}