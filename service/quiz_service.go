package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/persistence"
)

var (
	ErrValidation       = errors.New("validation failed")
	ErrResourceNotFound = errors.New("resource not found")
	ErrForbidden        = errors.New("forbidden")
)

const (
	nameEmptyError        = "name can not be null"
	invalidQuestionsError = "questions are invalid"
	teamNotFound          = "team not found"
	userNotFound          = "user not found"
	userNotInTeam         = "user not in team"
)

type QuizServiceInterface interface {
	CreateQuiz(request entity.Quiz) (dto.CreateQuizResponse, error)
}

type QuizService struct {
	teamRepo TeamRepositoryInterface
	userRepo UserRepositoryInterface
	quizRepo persistence.QuizRepositoryInterface
}

func NewQuizService() *QuizService {
	return &QuizService{
		teamRepo: persistence.NewTeamRepository(),
		userRepo: persistence.NewUserRepository(),
		quizRepo: persistence.NewQuizRepository(),
	}
}

func NewQuizServiceWithRepo(teamRepo TeamRepositoryInterface, userRepo UserRepositoryInterface, quizRepo persistence.QuizRepositoryInterface) *QuizService {
	return &QuizService{
		teamRepo: teamRepo,
		userRepo: userRepo,
		quizRepo: quizRepo,
	}
}

func (qs *QuizService) CreateQuiz(request entity.Quiz) (dto.CreateQuizResponse, error) {
	if request.QuizName == "" {
		return dto.CreateQuizResponse{}, fmt.Errorf("%w: %s", ErrValidation, nameEmptyError)
	}

	team, err := qs.teamRepo.GetTeamById(request.TeamID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return dto.CreateQuizResponse{}, fmt.Errorf("%w: %s", ErrResourceNotFound, teamNotFound)
		}
		return dto.CreateQuizResponse{}, err
	}

	user, err := qs.userRepo.GetByID(request.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return dto.CreateQuizResponse{}, fmt.Errorf("%w: %s", ErrResourceNotFound, userNotFound)
		}
		return dto.CreateQuizResponse{}, err
	}

	ok := 0

	if user.TeamsIds == nil {
		return dto.CreateQuizResponse{}, fmt.Errorf("%w: %s", ErrForbidden, teamNotFound)
	}

	teamIds := *user.TeamsIds

	for _, teamId := range teamIds {
		if teamId == team.Id {
			ok = 1
		}
	}

	if ok == 0 {
		return dto.CreateQuizResponse{}, fmt.Errorf("%w: %s", ErrForbidden, userNotInTeam)
	}

	for _, question := range request.Questions {
		if question.Question == "" || len(question.Options) == 0 || len(question.Answers) == 0 || question.Type == "" {
			return dto.CreateQuizResponse{}, fmt.Errorf("%w: %s", ErrValidation, invalidQuestionsError)
		}
	}

	id, err := generateID()
	if err != nil {
		return dto.CreateQuizResponse{}, err
	}
	request.ID = id

	err = qs.quizRepo.Create(request)
	if err != nil {
		return dto.CreateQuizResponse{}, err
	}

	return dto.NewCreateQuizResponse(id), nil
}
