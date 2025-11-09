package service

import (
	"errors"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/persistence"
	"github.com/SerbanEduard/ProiectColectivBackEnd/validator"
)

type TeamService struct {
	userRepository UserRepositoryInterface
	teamRepository TeamRepositoryInterface
}

type TeamRepositoryInterface interface {
	Create(team *entity.Team) error
	GetTeamById(id string) (*entity.Team, error)
	GetXTeamsByPrefix(prefix string, x int) ([]*entity.Team, error)
	GetTeamsByName(name string) ([]*entity.Team, error)
	GetAll() ([]*entity.Team, error)
	Update(team *entity.Team) error
	Delete(id string) error
}

func NewTeamService() *TeamService {
	return &TeamService{
		userRepository: persistence.NewUserRepository(),
		teamRepository: persistence.NewTeamRepository(),
	}
}

func NewTeamServiceWithRepo(UserRepositoryInterface UserRepositoryInterface, teamRepositoryInterface TeamRepositoryInterface) *TeamService {
	return &TeamService{
		userRepository: UserRepositoryInterface,
		teamRepository: teamRepositoryInterface,
	}
}

func (ts *TeamService) CreateTeam(request *dto.TeamRequest) (*dto.TeamResponse, error) {
	if err := validator.ValidateTeamRequest(request); err != nil {
		return nil, err
	}
	id, err := generateID()
	if err != nil {
		return nil, err
	}
	team := *entity.NewTeam(
		id,
		request.Name,
		request.Description,
		request.IsPublic,
		nil,
	)
	if err := ts.teamRepository.Create(&team); err != nil {
		return nil, err
	}
	return dto.NewTeamResponse(team.Name, team.Description, team.IsPublic), nil
}

func (ts *TeamService) AddUserToTeam(idUser string, idTeam string) error {
	user, err := ts.userRepository.GetByID(idUser)
	if err != nil {
		return err
	}
	team, err := ts.teamRepository.GetTeamById(idTeam)
	if err != nil {
		return err
	}
	for _, u := range team.UsersIds {
		if user.ID == u {
			return errors.New("User is already part of the team")
		}
	}
	team.UsersIds = append(team.UsersIds, idUser)
	if user.TeamsIds == nil {
		user.TeamsIds = &[]string{}
	}
	*user.TeamsIds = append(*user.TeamsIds, idTeam)

	if err := ts.userRepository.Update(user); err != nil {
		return err
	}
	if err := ts.teamRepository.Update(team); err != nil {
		return err
	}
	return nil
}

func (ts *TeamService) GetTeamById(id string) (*entity.Team, error) {
	return ts.teamRepository.GetTeamById(id)
}

func (ts *TeamService) GetXTeamsByPrefix(prefix string, x int) ([]*entity.Team, error) {
	return ts.teamRepository.GetXTeamsByPrefix(prefix, x)
}

func (ts *TeamService) GetTeamsByName(name string) ([]*entity.Team, error) {
	return ts.teamRepository.GetTeamsByName(name)
}

func (ts *TeamService) GetAll() ([]*entity.Team, error) {
	return ts.teamRepository.GetAll()
}

func (ts *TeamService) Update(team *entity.Team) error {
	return ts.teamRepository.Update(team)
}

// also deletes all references to the team in the Users' saved teams
func (ts *TeamService) Delete(id string) error {
	team, err := ts.teamRepository.GetTeamById(id)
	if err != nil {
		return err
	}
	for _, user := range team.UsersIds {
		user, err := ts.userRepository.GetByID(user)
		if err != nil {
			return err
		}
		if user.TeamsIds != nil {
			updatedTeams := removeString(*user.TeamsIds, team.Id)
			user.TeamsIds = &updatedTeams
		}
		if err := ts.userRepository.Update(user); err != nil {
			return err
		}
	}
	return ts.teamRepository.Delete(id)
}

func removeString(slice []string, value string) []string {
	result := []string{}
	for _, v := range slice {
		if v != value { // keep everything except the value
			result = append(result, v)
		}
	}
	return result
}
