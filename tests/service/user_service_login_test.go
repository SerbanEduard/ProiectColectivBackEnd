package service_test

import (
	"fmt"
	"testing"

	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Login_Success(t *testing.T) {
	mockRepo := new(tests.MockUserRepository)
	userService := service.NewUserServiceWithRepo(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	userInDB := tests.ValidUserWithPassword
	userInDB.Password = string(hashedPassword)

	request := &tests.ValidLoginRequest
	expectedResponse := &tests.ValidLoginResponse

	mockRepo.On("GetByEmail", "john@example.com").Return(&userInDB, nil)

	response, err := userService.Login(request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.FirstName, response.FirstName)
	assert.Equal(t, expectedResponse.LastName, response.LastName)
	assert.Equal(t, expectedResponse.Username, response.Username)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(tests.MockUserRepository)
	userService := service.NewUserServiceWithRepo(mockRepo)

	request := &tests.ValidLoginRequest

	mockRepo.On("GetByEmail", "john@example.com").Return(nil, fmt.Errorf("user not found"))

	response, err := userService.Login(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	mockRepo := new(tests.MockUserRepository)
	userService := service.NewUserServiceWithRepo(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("different_password"), bcrypt.DefaultCost)
	userInDB := tests.ValidUserWithPassword
	userInDB.Password = string(hashedPassword)

	request := &tests.ValidLoginRequest

	mockRepo.On("GetByEmail", "john@example.com").Return(&userInDB, nil)

	response, err := userService.Login(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}
