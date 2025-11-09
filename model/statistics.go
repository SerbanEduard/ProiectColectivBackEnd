package model

import "time"

type TimeSpentOnTeam struct {
	TeamId   string        `json:"teamId"`
	Duration time.Duration `json:"duration"`
}

type Statistics struct {
	TotalTimeSpentOnApp time.Duration     `json:"totalTimeSpentOnApp"`
	TimeSpentOnTeams    []TimeSpentOnTeam `json:"timeSpentOnTeams"`
}