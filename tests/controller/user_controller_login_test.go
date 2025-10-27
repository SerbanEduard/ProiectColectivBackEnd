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
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserController_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockUserService)
	userController := controller.NewUserControllerWithService(mockService)

	request := tests.ValidLoginRequest
	response := &tests.ValidLoginResponse

	mockService.On("Login", &request).Return(response, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	userController.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseBody dto.LoginUserResponse
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, response.ID, responseBody.ID)
	assert.Equal(t, response.Username, responseBody.Username)
	mockService.AssertExpectations(t)
}

func TestUserController_Login_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockUserService)
	userController := controller.NewUserControllerWithService(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	invalidJSON := []byte(`{"email": invalid}`)
	c.Request, _ = http.NewRequest("POST", "/users/login", bytes.NewBuffer(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	userController.Login(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Login")
}

func TestUserController_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockUserService)
	userController := controller.NewUserControllerWithService(mockService)

	request := tests.ValidLoginRequest

	mockService.On("Login", &request).Return(nil, fmt.Errorf("invalid credentials"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	userController.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var responseBody map[string]string
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, "invalid credentials", responseBody["error"])
	mockService.AssertExpectations(t)
}
