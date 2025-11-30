package entity

type TeamRequest struct {
	Id     string `json:"id"`
	UserID string `json:"userid"`
	TeamID string `json:"teamid"`
}

func NewTeamRequest(id, userId, teamId string) *TeamRequest {
	return &TeamRequest{
		Id:     id,
		UserID: userId,
		TeamID: teamId,
	}
}
