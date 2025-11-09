package dto

type TeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"ispublic"`
}

type TeamResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"ispublic"`
}

func NewTeamRequest(name, desc string, isPublic bool) *TeamRequest {
	return &TeamRequest{
		Name:        name,
		Description: desc,
		IsPublic:    isPublic,
	}
}

func NewTeamResponse(name, desc string, isPublic bool) *TeamResponse {
	return &TeamResponse{
		Name:        name,
		Description: desc,
		IsPublic:    isPublic,
	}
}
