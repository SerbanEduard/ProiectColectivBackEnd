package persistence

import (
	"context"

	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

type TeamRepository struct {
}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{}
}

func (tr *TeamRepository) Create(team *entity.Team) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef("teams/" + team.Id)
	return ref.Set(ctx, team)
}

func (tr *TeamRepository) GetTeamById(id string) (*entity.Team, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef("teams/" + id)

	var team entity.Team
	if err := ref.Get(ctx, &team); err != nil {
		return nil, err
	}
	return &team, nil
}

func (tr *TeamRepository) GetXTeamsByPrefix(prefix string, x int) ([]*entity.Team, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef("teams")

	query := ref.OrderByChild("name").
		StartAt(prefix).
		EndAt(prefix + "\uf8ff")

	results, err := query.GetOrdered(ctx)
	if err != nil {
		return nil, err
	}

	if len(results) > x {
		results = results[:x]
	}

	teams := make([]*entity.Team, 0, len(results))
	for _, r := range results {
		var team entity.Team
		if err := r.Unmarshal(&team); err != nil {
			return nil, err
		}
		teams = append(teams, &team)
	}

	return teams, nil
}

func (tr *TeamRepository) GetTeamsByName(name string) ([]*entity.Team, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef("teams")

	query := ref.OrderByChild("name").EqualTo(name)

	results, err := query.GetOrdered(ctx)
	if err != nil {
		return nil, err
	}

	teams := make([]*entity.Team, 0, len(results))
	for _, r := range results {
		var team entity.Team
		if err := r.Unmarshal(&team); err != nil {
			return nil, err
		}
		teams = append(teams, &team)
	}

	return teams, nil
}

func (tr *TeamRepository) GetAll() ([]*entity.Team, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef("teams")

	var teamsMap map[string]*entity.Team
	if err := ref.Get(ctx, &teamsMap); err != nil {
		return nil, err
	}

	teams := make([]*entity.Team, 0, len(teamsMap))
	for _, team := range teamsMap {
		teams = append(teams, team)
	}
	return teams, nil
}

func (tr *TeamRepository) Update(team *entity.Team) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef("teams/" + team.Id)
	return ref.Set(ctx, team)
}

func (tr *TeamRepository) Delete(id string) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef("teams/" + id)
	return ref.Delete(ctx)
}
