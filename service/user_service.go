package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/persistence"
	"github.com/SerbanEduard/ProiectColectivBackEnd/validator"
	"golang.org/x/crypto/bcrypt"
)

const (
	usernameAlreadyExistsError = "username already exists"
	emailAlreadyExistsError    = "email already exists"
	teamNotFoundError          = "team not found"
)

type UserService struct {
	userRepo UserRepositoryInterface
	teamRepo TeamRepositoryInterface
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: persistence.NewUserRepository(),
		teamRepo: persistence.NewTeamRepository(),
	}
}

func NewUserServiceWithRepo(userRepo interface{}, teamRepo interface{}) *UserService {
	return &UserService{
		userRepo: userRepo.(UserRepositoryInterface),
		teamRepo: teamRepo.(TeamRepositoryInterface),
	}
}

type UserRepositoryInterface interface {
	Create(user *entity.User) error
	GetByID(id string) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id string) error
	GetAll() ([]*entity.User, error)
}

func (us *UserService) SignUp(request *dto.SignUpUserRequest) (*dto.SignUpUserResponse, error) {
	if err := validator.ValidateSignUpRequest(request); err != nil {
		return nil, err
	}

	if _, err := us.userRepo.GetByUsername(request.Username); err == nil {
		return nil, fmt.Errorf(usernameAlreadyExistsError)
	}

	if _, err := us.userRepo.GetByEmail(request.Email); err == nil {
		return nil, fmt.Errorf(emailAlreadyExistsError)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	user := entity.NewUser(
		id,
		request.FirstName,
		request.LastName,
		request.Username,
		request.Email,
		string(hashedPassword),
		request.TopicsOfInterest,
	)

	if err := us.userRepo.Create(user); err != nil {
		return nil, err
	}

	return dto.NewSignUpUserResponse(user.FirstName, user.LastName, user.Username), nil
}

func (us *UserService) GetUserByID(id string) (*entity.User, error) {
	return us.userRepo.GetByID(id)
}

func (us *UserService) GetUserByEmail(email string) (*entity.User, error) {
	return us.userRepo.GetByEmail(email)
}

func (us *UserService) UpdateUser(user *entity.User) error {
	return us.userRepo.Update(user)
}

// also deletes all references to the user in the Teams' saved users
func (us *UserService) DeleteUser(id string) error {
	user, err := us.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	if user.TeamsIds == nil {
		return us.userRepo.Delete(id)
	}
	for _, teamId := range *user.TeamsIds {
		team, err := us.teamRepo.GetTeamById(teamId)
		if err != nil {
			return err
		}
		team.UsersIds = removeString(team.UsersIds, user.ID)
		if err := us.teamRepo.Update(team); err != nil {
			return err
		}
	}
	return us.userRepo.Delete(id)
}

func (us *UserService) GetAllUsers() ([]*entity.User, error) {
	return us.userRepo.GetAll()
}

func (us *UserService) UpdateUserStatistics(id string, timeSpentOnApp int64, timeSpentOnTeam model.TimeSpentOnTeam) (*entity.User, error) {
	user, err := us.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	_, err = us.teamRepo.GetTeamById(timeSpentOnTeam.TeamId)
	if err != nil {
		return nil, err
	}

	if user.Statistics == nil {
		user.Statistics = &model.Statistics{}
	}

	user.Statistics.TotalTimeSpentOnApp += timeSpentOnApp

	for i, teamTime := range user.Statistics.TimeSpentOnTeams {
		if teamTime.TeamId == timeSpentOnTeam.TeamId {
			user.Statistics.TimeSpentOnTeams[i].Duration += timeSpentOnTeam.Duration
			if err := us.userRepo.Update(user); err != nil {
				return nil, err
			}
			return user, nil
		}
	}

	user.Statistics.TimeSpentOnTeams = append(user.Statistics.TimeSpentOnTeams, timeSpentOnTeam)
	if err := us.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
