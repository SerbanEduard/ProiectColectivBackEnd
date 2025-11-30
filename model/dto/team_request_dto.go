package dto

import "github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"

type TeamRequestCreateDTO struct {
	UserID string `json:"userId" binding:"required"`
	TeamID string `json:"teamId" binding:"required"`
}

type TeamRequestItemDTO struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	TeamID string `json:"teamId"`
}

type TeamRequestsResponseDTO struct {
	Requests []TeamRequestItemDTO `json:"requests"`
}

func NewTeamRequestItemDTO(req *entity.TeamRequest) *TeamRequestItemDTO {
	return &TeamRequestItemDTO{
		ID:     req.Id,
		UserID: req.UserID,
		TeamID: req.TeamID,
	}
}

func NewTeamRequestsResponseDTO(reqs []*entity.TeamRequest) *TeamRequestsResponseDTO {
	items := make([]TeamRequestItemDTO, 0, len(reqs))
	for _, r := range reqs {
		items = append(items, *NewTeamRequestItemDTO(r))
	}
	return &TeamRequestsResponseDTO{Requests: items}
}
