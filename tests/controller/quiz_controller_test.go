package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	. "github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestQuizController_CreateQuiz_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockQuizService)
	qc := controller.NewQuizControllerWithService(mockService)

	request := entity.Quiz{
		QuizName: "Test Quiz",
		UserID:   TestUserID,
		TeamID:   TestTeamID,
		Questions: []entity.Question{
			{Question: "Q1", Options: []string{"a", "b"}, Answers: []string{"a"}, Type: "multiple_choice"},
		},
	}

	expectedResp := dto.NewCreateQuizResponse("created-id-1")
	mockService.On("CreateQuiz", request).Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/quizzes", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	qc.CreateQuiz(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.CreateQuizResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expectedResp.QuizID, resp.QuizID)

	mockService.AssertExpectations(t)
}

func TestQuizController_CreateQuiz_BadRequest_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockQuizService)
	qc := controller.NewQuizControllerWithService(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/quizzes", bytes.NewBuffer([]byte(`{invalid json}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	qc.CreateQuiz(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestQuizController_CreateQuiz_ServiceValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockQuizService)
	qc := controller.NewQuizControllerWithService(mockService)

	request := entity.Quiz{QuizName: ""}
	mockService.On("CreateQuiz", request).Return(dto.CreateQuizResponse{}, fmt.Errorf("%w: %s", service.ErrValidation, "name can not be null"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/quizzes", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	qc.CreateQuiz(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestQuizController_GetQuizWithAnswers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockQuizService)
	qc := controller.NewQuizControllerWithService(mockService)

	expectedQuiz := entity.Quiz{ID: "q-1", QuizName: "Quiz 1"}
	mockService.On("GetQuizWithAnswersById", "q-1").Return(expectedQuiz, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "q-1"}}

	qc.GetQuizWithAnswers(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp entity.Quiz
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expectedQuiz.ID, resp.ID)
	mockService.AssertExpectations(t)
}

func TestQuizController_GetQuizWithAnswers_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockQuizService)
	qc := controller.NewQuizControllerWithService(mockService)

	mockService.On("GetQuizWithAnswersById", "missing").Return(entity.Quiz{}, fmt.Errorf("%w: %s", service.ErrResourceNotFound, "quiz not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "missing"}}

	qc.GetQuizWithAnswers(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestQuizController_GetQuizWithoutAnswers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockQuizService)
	qc := controller.NewQuizControllerWithService(mockService)

	expectedResponse := dto.ReadQuizResponse{
		QuizID:    "q-1",
		QuizTitle: "Test Quiz",
		QuizQuestions: []dto.ReadQuizQuestionResponse{
			{QuestionID: "question-1", Question: "What is 2+2?", Options: []string{"3", "4", "5"}},
		},
	}

	mockService.On("GetQuizWithoutAnswersById", "q-1").Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "q-1"}}

	qc.GetQuizWithoutAnswers(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.ReadQuizResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expectedResponse.QuizID, resp.QuizID)
	assert.Equal(t, expectedResponse.QuizTitle, resp.QuizTitle)
	assert.Len(t, resp.QuizQuestions, 1)
	mockService.AssertExpectations(t)
}

func TestQuizController_GetQuizWithoutAnswers_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockQuizService)
	qc := controller.NewQuizControllerWithService(mockService)

	mockService.On("GetQuizWithoutAnswersById", "missing").Return(dto.ReadQuizResponse{}, fmt.Errorf("%w: %s", service.ErrResourceNotFound, "quiz not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "missing"}}

	qc.GetQuizWithoutAnswers(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestQuizController_GetQuizWithoutAnswers_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockQuizService)
	qc := controller.NewQuizControllerWithService(mockService)

	mockService.On("GetQuizWithoutAnswersById", "").Return(dto.ReadQuizResponse{}, fmt.Errorf("%w: %s", service.ErrValidation, "no id specified"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: ""}}

	qc.GetQuizWithoutAnswers(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestQuizController_SolveQuiz_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockQuizService)
	qc := controller.NewQuizControllerWithService(mockService)

	request := dto.SolveQuizRequest{
		QuizID: "q-1",
		Attempts: []dto.SolveQuestionRequest{
			{QuestionID: "question-1", Answer: []string{"4"}},
		},
	}

	expectedResponse := dto.SolveQuizResponse{
		IsCorrect: true,
		QuestionResponses: []dto.SolveQuestionResponse{
			{QuestionID: "question-1", IsCorrect: true, CorrectFields: []string{"4"}},
		},
	}

	mockService.On("SolveQuiz", request, TestUserID).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	claims := jwt.MapClaims{"sub": TestUserID}
	c.Set("userClaims", claims)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/quizzes/solve", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	qc.SolveQuiz(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.SolveQuizResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expectedResponse.IsCorrect, resp.IsCorrect)
	assert.Len(t, resp.QuestionResponses, 1)
	assert.Equal(t, expectedResponse.QuestionResponses[0].QuestionID, resp.QuestionResponses[0].QuestionID)
	mockService.AssertExpectations(t)
}

func TestQuizController_SolveQuiz_BadRequest_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockQuizService)
	qc := controller.NewQuizControllerWithService(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	claims := jwt.MapClaims{"sub": TestUserID}
	c.Set("userClaims", claims)

	c.Request, _ = http.NewRequest("POST", "/quizzes/solve", bytes.NewBuffer([]byte(`{invalid json}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	qc.SolveQuiz(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestQuizController_SolveQuiz_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockQuizService)
	qc := controller.NewQuizControllerWithService(mockService)

	request := dto.SolveQuizRequest{
		QuizID: "",
		Attempts: []dto.SolveQuestionRequest{
			{QuestionID: "question-1", Answer: []string{"4"}},
		},
	}

	mockService.On("SolveQuiz", request, TestUserID).Return(dto.SolveQuizResponse{}, fmt.Errorf("%w: %s", service.ErrValidation, "no id specified"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	claims := jwt.MapClaims{"sub": TestUserID}
	c.Set("userClaims", claims)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/quizzes/solve", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	qc.SolveQuiz(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestQuizController_SolveQuiz_QuizNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockQuizService)
	qc := controller.NewQuizControllerWithService(mockService)

	request := dto.SolveQuizRequest{
		QuizID: "missing-quiz",
		Attempts: []dto.SolveQuestionRequest{
			{QuestionID: "question-1", Answer: []string{"4"}},
		},
	}

	mockService.On("SolveQuiz", request, TestUserID).Return(dto.SolveQuizResponse{}, fmt.Errorf("%w: %s", service.ErrResourceNotFound, "quiz not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	claims := jwt.MapClaims{"sub": TestUserID}
	c.Set("userClaims", claims)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/quizzes/solve", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	qc.SolveQuiz(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestQuizController_SolveQuiz_IncorrectAnswers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockQuizService)
	qc := controller.NewQuizControllerWithService(mockService)

	request := dto.SolveQuizRequest{
		QuizID: "q-1",
		Attempts: []dto.SolveQuestionRequest{
			{QuestionID: "question-1", Answer: []string{"3"}},
		},
	}

	expectedResponse := dto.SolveQuizResponse{
		IsCorrect: false,
		QuestionResponses: []dto.SolveQuestionResponse{
			{QuestionID: "question-1", IsCorrect: false, CorrectFields: []string{"4"}},
		},
	}

	mockService.On("SolveQuiz", request, TestUserID).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	claims := jwt.MapClaims{"sub": TestUserID}
	c.Set("userClaims", claims)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/quizzes/solve", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	qc.SolveQuiz(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.SolveQuizResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp.IsCorrect)
	assert.Len(t, resp.QuestionResponses, 1)
	assert.False(t, resp.QuestionResponses[0].IsCorrect)
	mockService.AssertExpectations(t)
}
