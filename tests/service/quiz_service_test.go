package service_test

import (
	"errors"
	"testing"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	. "github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const MockQuizID = "mock-quiz-id"
const QuizNotFoundErr = "quiz not found"
const TeamNotFoundErr = "team not found"
const UserNotFoundErr = "user not found"

func getValidQuizRequestEntity() entity.Quiz {
	return entity.Quiz{
		QuizName: "Test Quiz Name",
		TeamID:   TestTeamID,
		UserID:   TestUserID,
		Questions: []entity.Question{
			{
				Question: "What is 2+2?",
				Options:  []string{"3", "4", "5"},
				Answers:  []string{"4"},
				Type:     "single",
			},
		},
	}
}

func TestQuizService_CreateQuiz_Success(t *testing.T) {
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	mockQuizRepo := new(tests.MockQuizRepository)

	quizService := service.NewQuizServiceWithRepo(mockTeamRepo, mockUserRepo, mockQuizRepo)
	request := getValidQuizRequestEntity()

	team := &entity.Team{Id: TestTeamID}
	userTeams := []string{TestTeamID, "other-team-id"}
	user := &entity.User{ID: TestUserID, TeamsIds: &userTeams}

	mockTeamRepo.On("GetTeamById", TestTeamID).Return(team, nil).Once()
	mockUserRepo.On("GetByID", TestUserID).Return(user, nil).Once()
	mockQuizRepo.On("Create", mock.MatchedBy(func(q entity.Quiz) bool {
		return q.QuizName == request.QuizName && q.ID != ""
	})).Return(nil).Once()

	response, err := quizService.CreateQuiz(request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.QuizID)

	mockTeamRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockQuizRepo.AssertExpectations(t)
}

func TestQuizService_CreateQuiz_EmptyQuizName_ValidationFail(t *testing.T) {
	mockService := service.NewQuizServiceWithRepo(nil, nil, nil)
	request := getValidQuizRequestEntity()
	request.QuizName = ""

	_, err := mockService.CreateQuiz(request)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrValidation))
	assert.Contains(t, err.Error(), "name can not be null")
}

func TestQuizService_CreateQuiz_TeamNotFound(t *testing.T) {
	mockTeamRepo := new(tests.MockTeamRepository)
	quizService := service.NewQuizServiceWithRepo(mockTeamRepo, nil, nil)
	request := getValidQuizRequestEntity()

	mockTeamRepo.On("GetTeamById", TestTeamID).Return(nil, errors.New("db error: not found")).Once()

	_, err := quizService.CreateQuiz(request)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrResourceNotFound))
	assert.Contains(t, err.Error(), TeamNotFoundErr)
	mockTeamRepo.AssertExpectations(t)
}

func TestQuizService_CreateQuiz_UserNotFound(t *testing.T) {
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	quizService := service.NewQuizServiceWithRepo(mockTeamRepo, mockUserRepo, nil)
	request := getValidQuizRequestEntity()
	team := &entity.Team{Id: TestTeamID}

	mockTeamRepo.On("GetTeamById", TestTeamID).Return(team, nil).Once()
	mockUserRepo.On("GetByID", TestUserID).Return(nil, errors.New("db error: document not found")).Once()

	_, err := quizService.CreateQuiz(request)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrResourceNotFound))
	assert.Contains(t, err.Error(), UserNotFoundErr)
	mockTeamRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestQuizService_CreateQuiz_UserNotAMember_Forbidden(t *testing.T) {
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	quizService := service.NewQuizServiceWithRepo(mockTeamRepo, mockUserRepo, nil)
	request := getValidQuizRequestEntity()
	team := &entity.Team{Id: TestTeamID}

	otherTeams := []string{"team-001", "team-002"}
	user := &entity.User{ID: TestUserID, TeamsIds: &otherTeams}
	mockTeamRepo.On("GetTeamById", TestTeamID).Return(team, nil).Once()
	mockUserRepo.On("GetByID", TestUserID).Return(user, nil).Once()

	_, err := quizService.CreateQuiz(request)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrForbidden))
	assert.Contains(t, err.Error(), "user not in team")
	mockTeamRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestQuizService_CreateQuiz_InvalidQuestionFormat_ValidationFail(t *testing.T) {
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	quizService := service.NewQuizServiceWithRepo(mockTeamRepo, mockUserRepo, nil)
	request := getValidQuizRequestEntity()
	team := &entity.Team{Id: TestTeamID}
	userTeams := []string{TestTeamID}
	user := &entity.User{ID: TestUserID, TeamsIds: &userTeams}

	request.Questions[0].Options = []string{}
	mockTeamRepo.On("GetTeamById", TestTeamID).Return(team, nil).Once()
	mockUserRepo.On("GetByID", TestUserID).Return(user, nil).Once()

	_, err := quizService.CreateQuiz(request)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrValidation))
	assert.Contains(t, err.Error(), "questions are invalid")
	mockTeamRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestQuizService_GetQuizWithAnswersById_Success(t *testing.T) {
	mockQuizRepo := new(tests.MockQuizRepository)
	quizService := service.NewQuizServiceWithRepo(nil, nil, mockQuizRepo)

	expectedQuiz := getValidQuizRequestEntity()
	expectedQuiz.ID = MockQuizID

	mockQuizRepo.On("GetById", MockQuizID).Return(expectedQuiz, nil).Once()

	resultQuiz, err := quizService.GetQuizWithAnswersById(MockQuizID)
	assert.NoError(t, err)
	assert.Equal(t, expectedQuiz, resultQuiz)
	mockQuizRepo.AssertExpectations(t)
}

func TestQuizService_GetQuizWithAnswersById_EmptyID_ValidationFail(t *testing.T) {
	mockService := service.NewQuizServiceWithRepo(nil, nil, nil)

	_, err := mockService.GetQuizWithAnswersById("")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrValidation))
	assert.Contains(t, err.Error(), "no id specified")
}

func TestQuizService_GetQuizWithAnswersById_NotFound(t *testing.T) {
	mockQuizRepo := new(tests.MockQuizRepository)
	quizService := service.NewQuizServiceWithRepo(nil, nil, mockQuizRepo)

	mockQuizRepo.On("GetById", MockQuizID).Return(entity.Quiz{}, errors.New("db error: quiz not found")).Once()

	_, err := quizService.GetQuizWithAnswersById(MockQuizID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrResourceNotFound))
	assert.Contains(t, err.Error(), QuizNotFoundErr)
	mockQuizRepo.AssertExpectations(t)
}

func TestQuizService_GetQuizWithoutAnswersById_Success(t *testing.T) {
	mockQuizRepo := new(tests.MockQuizRepository)
	quizService := service.NewQuizServiceWithRepo(nil, nil, mockQuizRepo)

	quiz := getValidQuizRequestEntity()
	quiz.ID = MockQuizID
	quiz.Questions[0].ID = "question-1"

	mockQuizRepo.On("GetById", MockQuizID).Return(quiz, nil).Once()

	result, err := quizService.GetQuizWithoutAnswersById(MockQuizID)

	assert.NoError(t, err)
	assert.Equal(t, MockQuizID, result.QuizID)
	assert.Equal(t, quiz.QuizName, result.QuizTitle)
	assert.Len(t, result.QuizQuestions, 1)
	assert.Equal(t, "question-1", result.QuizQuestions[0].QuestionID)
	assert.Equal(t, "What is 2+2?", result.QuizQuestions[0].Question)
	assert.Equal(t, []string{"3", "4", "5"}, result.QuizQuestions[0].Options)
	mockQuizRepo.AssertExpectations(t)
}

func TestQuizService_GetQuizWithoutAnswersById_EmptyID_ValidationFail(t *testing.T) {
	quizService := service.NewQuizServiceWithRepo(nil, nil, nil)

	_, err := quizService.GetQuizWithoutAnswersById("")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrValidation))
	assert.Contains(t, err.Error(), "no id specified")
}

func TestQuizService_GetQuizWithoutAnswersById_NotFound(t *testing.T) {
	mockQuizRepo := new(tests.MockQuizRepository)
	quizService := service.NewQuizServiceWithRepo(nil, nil, mockQuizRepo)

	mockQuizRepo.On("GetById", MockQuizID).Return(entity.Quiz{}, errors.New("db error: quiz not found")).Once()

	_, err := quizService.GetQuizWithoutAnswersById(MockQuizID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrResourceNotFound))
	assert.Contains(t, err.Error(), QuizNotFoundErr)
	mockQuizRepo.AssertExpectations(t)
}

func TestQuizService_SolveQuiz_Success_AllCorrect(t *testing.T) {
	mockQuizRepo := new(tests.MockQuizRepository)
	mockUserRepo := new(tests.MockUserRepository)
	quizService := service.NewQuizServiceWithRepo(nil, mockUserRepo, mockQuizRepo)

	quiz := getValidQuizRequestEntity()
	quiz.ID = MockQuizID
	quiz.Questions[0].ID = "question-1"

	userTeams := []string{TestTeamID}
	user := &entity.User{ID: TestUserID, TeamsIds: &userTeams}

	solveRequest := dto.SolveQuizRequest{
		QuizID: MockQuizID,
		Attempts: []dto.SolveQuestionRequest{
			{QuestionID: "question-1", Answer: []string{"4"}},
		},
	}

	mockQuizRepo.On("GetById", MockQuizID).Return(quiz, nil).Once()
	mockUserRepo.On("GetByID", TestUserID).Return(user, nil).Once()

	result, err := quizService.SolveQuiz(solveRequest, TestUserID)

	assert.NoError(t, err)
	assert.True(t, result.IsCorrect)
	assert.Len(t, result.QuestionResponses, 1)
	assert.Equal(t, "question-1", result.QuestionResponses[0].QuestionID)
	assert.True(t, result.QuestionResponses[0].IsCorrect)
	assert.Equal(t, []string{"4"}, result.QuestionResponses[0].CorrectFields)
	mockQuizRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestQuizService_SolveQuiz_Success_SomeIncorrect(t *testing.T) {
	mockQuizRepo := new(tests.MockQuizRepository)
	mockUserRepo := new(tests.MockUserRepository)
	quizService := service.NewQuizServiceWithRepo(nil, mockUserRepo, mockQuizRepo)

	quiz := getValidQuizRequestEntity()
	quiz.ID = MockQuizID
	quiz.Questions[0].ID = "question-1"

	userTeams := []string{TestTeamID}
	user := &entity.User{ID: TestUserID, TeamsIds: &userTeams}

	solveRequest := dto.SolveQuizRequest{
		QuizID: MockQuizID,
		Attempts: []dto.SolveQuestionRequest{
			{QuestionID: "question-1", Answer: []string{"3"}},
		},
	}

	mockQuizRepo.On("GetById", MockQuizID).Return(quiz, nil).Once()
	mockUserRepo.On("GetByID", TestUserID).Return(user, nil).Once()

	result, err := quizService.SolveQuiz(solveRequest, TestUserID)

	assert.NoError(t, err)
	assert.False(t, result.IsCorrect)
	assert.Len(t, result.QuestionResponses, 1)
	assert.Equal(t, "question-1", result.QuestionResponses[0].QuestionID)
	assert.False(t, result.QuestionResponses[0].IsCorrect)
	assert.Equal(t, []string{"4"}, result.QuestionResponses[0].CorrectFields)
	mockQuizRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestQuizService_SolveQuiz_EmptyQuizID_ValidationFail(t *testing.T) {
	quizService := service.NewQuizServiceWithRepo(nil, nil, nil)

	solveRequest := dto.SolveQuizRequest{
		QuizID: "",
		Attempts: []dto.SolveQuestionRequest{
			{QuestionID: "question-1", Answer: []string{"4"}},
		},
	}

	_, err := quizService.SolveQuiz(solveRequest, TestUserID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrValidation))
	assert.Contains(t, err.Error(), "no id specified")
}

func TestQuizService_SolveQuiz_QuizNotFound(t *testing.T) {
	mockQuizRepo := new(tests.MockQuizRepository)
	quizService := service.NewQuizServiceWithRepo(nil, nil, mockQuizRepo)

	solveRequest := dto.SolveQuizRequest{
		QuizID: MockQuizID,
		Attempts: []dto.SolveQuestionRequest{
			{QuestionID: "question-1", Answer: []string{"4"}},
		},
	}

	mockQuizRepo.On("GetById", MockQuizID).Return(entity.Quiz{}, errors.New("db error: quiz not found")).Once()

	_, err := quizService.SolveQuiz(solveRequest, TestUserID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrResourceNotFound))
	assert.Contains(t, err.Error(), QuizNotFoundErr)
	mockQuizRepo.AssertExpectations(t)
}

func TestQuizService_SolveQuiz_MultipleQuestions_MixedResults(t *testing.T) {
	mockQuizRepo := new(tests.MockQuizRepository)
	mockUserRepo := new(tests.MockUserRepository)
	quizService := service.NewQuizServiceWithRepo(nil, mockUserRepo, mockQuizRepo)

	quiz := entity.Quiz{
		ID:       MockQuizID,
		QuizName: "Multi Question Quiz",
		TeamID:   TestTeamID,
		UserID:   TestUserID,
		Questions: []entity.Question{
			{ID: "q1", Question: "What is 2+2?", Options: []string{"3", "4", "5"}, Answers: []string{"4"}, Type: "single"},
			{ID: "q2", Question: "What is 3+3?", Options: []string{"6", "7", "8"}, Answers: []string{"6"}, Type: "single"},
		},
	}

	userTeams := []string{TestTeamID}
	user := &entity.User{ID: TestUserID, TeamsIds: &userTeams}

	solveRequest := dto.SolveQuizRequest{
		QuizID: MockQuizID,
		Attempts: []dto.SolveQuestionRequest{
			{QuestionID: "q1", Answer: []string{"4"}},
			{QuestionID: "q2", Answer: []string{"7"}},
		},
	}

	mockQuizRepo.On("GetById", MockQuizID).Return(quiz, nil).Once()
	mockUserRepo.On("GetByID", TestUserID).Return(user, nil).Once()

	result, err := quizService.SolveQuiz(solveRequest, TestUserID)

	assert.NoError(t, err)
	assert.False(t, result.IsCorrect)
	assert.Len(t, result.QuestionResponses, 2)
	assert.Equal(t, "q1", result.QuestionResponses[0].QuestionID)
	assert.True(t, result.QuestionResponses[0].IsCorrect)
	assert.Equal(t, []string{"4"}, result.QuestionResponses[0].CorrectFields)

	assert.Equal(t, "q2", result.QuestionResponses[1].QuestionID)
	assert.False(t, result.QuestionResponses[1].IsCorrect)
	assert.Equal(t, []string{"6"}, result.QuestionResponses[1].CorrectFields)

	mockQuizRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestQuizService_SolveQuiz_UserNotInTeam_Forbidden(t *testing.T) {
	mockQuizRepo := new(tests.MockQuizRepository)
	mockUserRepo := new(tests.MockUserRepository)
	quizService := service.NewQuizServiceWithRepo(nil, mockUserRepo, mockQuizRepo)

	quiz := getValidQuizRequestEntity()
	quiz.ID = MockQuizID
	quiz.Questions[0].ID = "question-1"

	otherTeams := []string{"other-team-1", "other-team-2"}
	user := &entity.User{ID: TestUserID, TeamsIds: &otherTeams}

	solveRequest := dto.SolveQuizRequest{
		QuizID: MockQuizID,
		Attempts: []dto.SolveQuestionRequest{
			{QuestionID: "question-1", Answer: []string{"4"}},
		},
	}

	mockQuizRepo.On("GetById", MockQuizID).Return(quiz, nil).Once()
	mockUserRepo.On("GetByID", TestUserID).Return(user, nil).Once()

	_, err := quizService.SolveQuiz(solveRequest, TestUserID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrForbidden))
	assert.Contains(t, err.Error(), "user not in team")
	mockQuizRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
