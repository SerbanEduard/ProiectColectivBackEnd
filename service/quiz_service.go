package service

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/SerbanEduard/ProiectColectivBackEnd/mappers"
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
	quizIdEmpty           = "no id specified"
	quizNotFound          = "quiz not found"
	NotFoundError         = "not found"
)

type QuizServiceInterface interface {
	CreateQuiz(request entity.Quiz) (dto.CreateQuizResponse, error)
	GetQuizWithAnswersById(id string) (entity.Quiz, error)
	GetQuizWithoutAnswersById(id string) (dto.ReadQuizResponse, error)
	SolveQuiz(request dto.SolveQuizRequest, userId string) (dto.SolveQuizResponse, error)
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

func (qs *QuizService) isUserInTeam(userId string, teamId string) (bool, error) {
	user, err := qs.userRepo.GetByID(userId)
	if err != nil {
		if strings.Contains(err.Error(), NotFoundError) {
			return false, fmt.Errorf("%w: %s", ErrResourceNotFound, userNotFound)
		}
		return false, err
	}

	ok := 0

	if user.TeamsIds == nil {
		return false, fmt.Errorf("%w: %s", ErrForbidden, teamNotFound)
	}

	teamIds := *user.TeamsIds

	for _, team := range teamIds {
		if team == teamId {
			ok = 1
		}
	}

	if ok == 0 {
		return false, fmt.Errorf("%w: %s", ErrForbidden, userNotInTeam)
	}

	return true, nil
}

func (qs *QuizService) CreateQuiz(request entity.Quiz) (dto.CreateQuizResponse, error) {
	if request.QuizName == "" {
		return dto.CreateQuizResponse{}, fmt.Errorf("%w: %s", ErrValidation, nameEmptyError)
	}

	team, err := qs.teamRepo.GetTeamById(request.TeamID)
	if err != nil {
		if strings.Contains(err.Error(), NotFoundError) {
			return dto.CreateQuizResponse{}, fmt.Errorf("%w: %s", ErrResourceNotFound, teamNotFound)
		}
		return dto.CreateQuizResponse{}, err
	}

	if isPartOf, err := qs.isUserInTeam(request.UserID, team.Id); isPartOf == false {
		return dto.CreateQuizResponse{}, err
	}

	for i, question := range request.Questions {
		if question.Question == "" || len(question.Options) == 0 || len(question.Answers) == 0 || question.Type == "" {
			return dto.CreateQuizResponse{}, fmt.Errorf("%w: %s", ErrValidation, invalidQuestionsError)
		}

		questionID, err := generateID()
		if err != nil {
			return dto.CreateQuizResponse{}, err
		}
		request.Questions[i].ID = questionID
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

func (qs *QuizService) GetQuizWithAnswersById(id string) (entity.Quiz, error) {
	if id == "" {
		return entity.Quiz{}, fmt.Errorf("%w: %s", ErrValidation, quizIdEmpty)
	}

	quiz, err := qs.quizRepo.GetById(id)
	if err != nil {
		if strings.Contains(err.Error(), NotFoundError) {
			return entity.Quiz{}, fmt.Errorf("%w: %s", ErrResourceNotFound, quizNotFound)
		}
		return entity.Quiz{}, err
	}

	return quiz, nil
}

func (qs *QuizService) GetQuizWithoutAnswersById(id string) (dto.ReadQuizResponse, error) {
	if id == "" {
		return dto.ReadQuizResponse{}, fmt.Errorf("%w: %s", ErrValidation, quizIdEmpty)
	}

	quiz, err := qs.quizRepo.GetById(id)
	if err != nil {
		if strings.Contains(err.Error(), NotFoundError) {
			return dto.ReadQuizResponse{}, fmt.Errorf("%w: %s", ErrResourceNotFound, quizNotFound)
		}
		return dto.ReadQuizResponse{}, err
	}

	quizWithoutAnswers := mappers.MapDomainToReadDTO(quiz)
	return quizWithoutAnswers, nil
}

func (qs *QuizService) SolveQuiz(request dto.SolveQuizRequest, userId string) (dto.SolveQuizResponse, error) {
	if request.QuizID == "" {
		return dto.SolveQuizResponse{}, fmt.Errorf("%w: %s", ErrValidation, quizIdEmpty)
	}

	questionsSubmitted := request.Attempts
	sort.Slice(questionsSubmitted, func(i, j int) bool {
		return questionsSubmitted[i].QuestionID < questionsSubmitted[j].QuestionID
	})
	quiz, err := qs.quizRepo.GetById(request.QuizID)
	if err != nil {
		if strings.Contains(err.Error(), NotFoundError) {
			return dto.SolveQuizResponse{}, fmt.Errorf("%w: %s", ErrResourceNotFound, quizNotFound)
		}
		return dto.SolveQuizResponse{}, err
	}

	if isPartOf, err := qs.isUserInTeam(userId, quiz.TeamID); isPartOf == false {
		return dto.SolveQuizResponse{}, err
	}

	questions := quiz.Questions
	sort.Slice(questions, func(i, j int) bool {
		return questions[i].ID < questions[j].ID
	})

	if len(questionsSubmitted) != len(questions) {
		return dto.SolveQuizResponse{}, fmt.Errorf("%w: %s", ErrValidation, invalidQuestionsError)
	}

	for i, question := range questions {
		submitted := questionsSubmitted[i]
		if question.ID != submitted.QuestionID {
			return dto.SolveQuizResponse{}, fmt.Errorf("%w: %s", ErrValidation, invalidQuestionsError)
		}
	}

	allCorrect := true
	questionResponses := make([]dto.SolveQuestionResponse, len(questions))

	for i, question := range questions {
		submitted := questionsSubmitted[i]
		correctFields := question.Answers
		sort.Slice(correctFields, func(i, j int) bool { return correctFields[i] < correctFields[j] })
		submittedFields := submitted.Answer
		sort.Slice(submittedFields, func(i, j int) bool { return submittedFields[i] < submittedFields[j] })
		isCorrect := slices.Equal(correctFields, submittedFields)
		if !isCorrect {
			allCorrect = false
		}
		questionResponses[i] = dto.NewSolveQuestionResponse(question.ID, isCorrect, correctFields)
	}

	return dto.SolveQuizResponse{
		IsCorrect:         allCorrect,
		QuestionResponses: questionResponses,
	}, nil
}
