package service_test

import (
	"fmt"
	"testing"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_SignUp_Success(t *testing.T) {
	mockRepo := new(tests.MockUserRepository)
	mockRepoTeams := new(tests.MockTeamsRepository)
	userService := service.NewUserServiceWithRepo(mockRepo, mockRepoTeams)

	request := &tests.ValidSignUpRequest

	mockRepo.On("GetByUsername", "johndoe").Return(nil, fmt.Errorf("user not found"))
	mockRepo.On("GetByEmail", "john@example.com").Return(nil, fmt.Errorf("user not found"))
	mockRepo.On("Create", mock.MatchedBy(func(user *entity.User) bool {
		return user.FirstName == "John" &&
			user.LastName == "Doe" &&
			user.Username == "johndoe" &&
			user.Email == "john@example.com" &&
			len(user.TopicsOfInterest) == 1 &&
			user.TopicsOfInterest[0] == model.Programming &&
			user.ID != "" &&
			user.Password != "password123"
	})).Return(nil)

	response, err := userService.SignUp(request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "John", response.FirstName)
	assert.Equal(t, "Doe", response.LastName)
	mockRepo.AssertExpectations(t)
}

func TestUserService_SignUp_UsernameExists(t *testing.T) {
	mockRepo := new(tests.MockUserRepository)
	mockRepoTeams := new(tests.MockTeamsRepository)
	userService := service.NewUserServiceWithRepo(mockRepo, mockRepoTeams)

	request := &tests.ExistingUsernameRequest

	existingUser := &tests.ExistingUser
	mockRepo.On("GetByUsername", "existinguser").Return(existingUser, nil)

	response, err := userService.SignUp(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "username already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserService_SignUp_EmailExists(t *testing.T) {
	mockRepo := new(tests.MockUserRepository)
	mockRepoTeams := new(tests.MockTeamsRepository)
	userService := service.NewUserServiceWithRepo(mockRepo, mockRepoTeams)

	request := &tests.ExistingEmailRequest

	mockRepo.On("GetByUsername", "johndoe").Return(nil, fmt.Errorf("user not found"))
	mockRepo.On("GetByEmail", "existing@example.com").Return(&tests.ExistingUser, nil)

	response, err := userService.SignUp(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "email already exists", err.Error())
	mockRepo.AssertExpectations(t)
}
