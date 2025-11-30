package persistence

import (
	"context"
	"errors"

	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

const (
	teamRequestsCollection = "teamRequests"
	requestNotFound        = "team request not found"
)

type TeamRequestRepository struct{}

func NewTeamRequestRepository() *TeamRequestRepository {
	return &TeamRequestRepository{}
}

func (tr *TeamRequestRepository) Create(req *entity.TeamRequest) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(teamRequestsCollection + "/" + req.Id)
	return ref.Set(ctx, req)
}

func (tr *TeamRequestRepository) GetById(id string) (*entity.TeamRequest, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(teamRequestsCollection + "/" + id)

	var req entity.TeamRequest
	if err := ref.Get(ctx, &req); err != nil {
		return nil, err
	}
	if req.Id == "" {
		return nil, errors.New(requestNotFound)
	}
	return &req, nil
}

func (tr *TeamRequestRepository) Delete(id string) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(teamRequestsCollection + "/" + id)
	return ref.Delete(ctx)
}

func (tr *TeamRequestRepository) GetAll() ([]*entity.TeamRequest, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(teamRequestsCollection)

	var reqsMap map[string]*entity.TeamRequest
	if err := ref.Get(ctx, &reqsMap); err != nil {
		return nil, err
	}
	reqs := make([]*entity.TeamRequest, 0, len(reqsMap))
	for _, r := range reqsMap {
		reqs = append(reqs, r)
	}
	return reqs, nil
}

func (tr *TeamRequestRepository) GetByUserId(userId string) ([]*entity.TeamRequest, error) {
	all, err := tr.GetAll()
	if err != nil {
		return nil, err
	}

	filtered := make([]*entity.TeamRequest, 0)
	for _, r := range all {
		if r.UserID == userId {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}
