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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserController_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockUserService)
	userController := controller.NewUserControllerWithService(mockService)

	// prepare hashed password same as service stores
	hashed, _ := bcrypt.GenerateFromPassword([]byte(tests.TestPassword), bcrypt.DefaultCost)

	topics := &[]model.TopicOfInterest{model.Programming}
	user := &entity.User{
		ID:               tests.TestUserID,
		Username:         tests.TestUsername,
		Email:            tests.TestEmail,
		Password:         string(hashed),
		TopicsOfInterest: topics,
	}

	// mock service behavior
	mockService.On("GetUserByEmail", tests.TestEmail).Return(user, nil)

	// build request
	loginReq := dto.LoginRequest{Email: tests.TestEmail, Password: tests.TestPassword}
	jsonData, _ := json.Marshal(loginReq)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	userController.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.Equal(t, tests.TestUserID, resp.User.Id)
	assert.Equal(t, tests.TestUsername, resp.User.Username)
	assert.Equal(t, tests.TestEmail, resp.User.Email)

	mockService.AssertExpectations(t)
}

func TestUserController_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockUserService)
	userController := controller.NewUserControllerWithService(mockService)

	// service returns error (user not found)
	mockService.On("GetUserByEmail", tests.TestEmail).Return(nil, fmt.Errorf("not found"))

	loginReq := dto.LoginRequest{Email: tests.TestEmail, Password: "wrongpass"}
	jsonData, _ := json.Marshal(loginReq)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	userController.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockService.AssertExpectations(t)
}
