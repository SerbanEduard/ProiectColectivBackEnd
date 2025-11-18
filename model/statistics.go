package model

type TimeSpentOnTeam struct {
	TeamId   string `json:"teamId"`
	Duration int64  `json:"duration" example:"3600000" description:"Duration in milliseconds"`
}

type Statistics struct {
	TotalTimeSpentOnApp int64             `json:"totalTimeSpentOnApp" example:"7200000" description:"Total time spent on app in milliseconds"`
	TimeSpentOnTeams    []TimeSpentOnTeam `json:"timeSpentOnTeams"`
}
