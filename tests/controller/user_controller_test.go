package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	. "github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	TestDuration150Min = int64(9000000) // 150 minutes in milliseconds
)

func TestUserController_SignUp_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockUserService)
	userController := controller.NewUserControllerWithService(mockService)

	request := tests.ValidSignUpRequest

	response := &tests.ValidSignUpResponse

	mockService.On("SignUp", &request).Return(response, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest(HTTPMethodPOST, PathUsersSignup, bytes.NewBuffer(jsonData))
	c.Request.Header.Set(ContentTypeJSON, ContentTypeJSON)

	userController.SignUp(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestUserController_SignUp_UsernameExists(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockUserService)
	userController := controller.NewUserControllerWithService(mockService)

	request := tests.ExistingUsernameRequest

	mockService.On("SignUp", &request).Return(nil, fmt.Errorf(ErrUsernameExists))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest(HTTPMethodPOST, PathUsersSignup, bytes.NewBuffer(jsonData))
	c.Request.Header.Set(ContentTypeJSON, ContentTypeJSON)

	userController.SignUp(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var responseBody map[string]string
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, ErrUsernameExists, responseBody[JSONKeyError])

	mockService.AssertExpectations(t)
}

func TestUserController_SignUp_EmailExists(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockUserService)
	userController := controller.NewUserControllerWithService(mockService)

	request := tests.ExistingEmailRequest

	mockService.On("SignUp", &request).Return(nil, fmt.Errorf("email already exists"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	userController.SignUp(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var responseBody map[string]string
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, "email already exists", responseBody["error"])

	mockService.AssertExpectations(t)
}

func TestUserController_UpdateUserStatistics_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockUserService)
	userController := controller.NewUserControllerWithService(mockService)

	request := tests.ValidUpdateStatisticsRequest

	mockService.On("UpdateUserStatistics", TestUserID, TestDurationApp, tests.ValidTimeSpentOnTeam).Return(&entity.User{ID: TestUserID, Statistics: &model.Statistics{}}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest(HTTPMethodPUT, PathUserStatistics, bytes.NewBuffer(jsonData))
	c.Request.Header.Set(ContentTypeJSON, ContentTypeJSON)
	c.Params = []gin.Param{{Key: ParamKeyID, Value: TestUserID}}

	userController.UpdateUserStatistics(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseBody dto.UpdateStatisticsResponse
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, TestUserID, responseBody.UserId)

	mockService.AssertExpectations(t)
}
