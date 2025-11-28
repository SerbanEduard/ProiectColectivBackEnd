package dto

import "github.com/SerbanEduard/ProiectColectivBackEnd/model"

type UpdateStatisticsRequest struct {
	TimeSpentOnApp  int64  `json:"timeSpentOnApp" example:"1800000" description:"Time spent on app in milliseconds"`
	TeamId          string `json:"teamId" example:"team123"`
	TimeSpentOnTeam int64  `json:"timeSpentOnTeam" example:"900000" description:"Time spent on team in milliseconds"`
}

type StatisticsResponse struct {
	UserId              string                  `json:"userId"`
	TotalTimeSpentOnApp int64                   `json:"totalTimeSpentOnApp" example:"7200000" description:"Total time spent on app in milliseconds"`
	TimeSpentOnTeams    []model.TimeSpentOnTeam `json:"timeSpentOnTeams"`
}

func NewStatisticsResponse(userId string, statistics *model.Statistics) *StatisticsResponse {
	return &StatisticsResponse{
		UserId:              userId,
		TotalTimeSpentOnApp: statistics.TotalTimeSpentOnApp,
		TimeSpentOnTeams:    statistics.TimeSpentOnTeams,
	}
}
