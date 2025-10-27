package service_test

import (
	"fmt"
	"testing"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFriendRequestService_SendFriendRequest_Success(t *testing.T) {
	mockRepo := new(tests.MockFriendRequestRepository)
	mockUserService := new(tests.MockUserService)

	service := service.NewFriendRequestService()
	service.SetFriendRequestRepo(mockRepo)
	service.SetUserService(mockUserService)

	mockUserService.On("GetUserByID", "user1").Return(&entity.User{ID: "user1"}, nil)
	mockUserService.On("GetUserByID", "user2").Return(&entity.User{ID: "user2"}, nil)
	mockRepo.On("GetByUsers", "user1", "user2").Return(nil, fmt.Errorf("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*entity.FriendRequest")).Return(nil)

	err := service.SendFriendRequest("user1", "user2")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

func TestFriendRequestService_RespondToRequest_Accept(t *testing.T) {
	mockRepo := new(tests.MockFriendRequestRepository)
	service := service.NewFriendRequestService()
	service.SetFriendRequestRepo(mockRepo)

	request := tests.ValidFriendRequest
	mockRepo.On("GetByUsers", "user1", "user2").Return(&request, nil)
	mockRepo.On("Update", mock.MatchedBy(func(r *entity.FriendRequest) bool {
		return r.Status == entity.ACCEPTED
	})).Return(nil)

	err := service.RespondToFriendRequest("user1", "user2", true)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestFriendRequestService_SendFriendRequest_InvalidSenderID(t *testing.T) {
	mockRepo := new(tests.MockFriendRequestRepository)
	mockUserService := new(tests.MockUserService)
	service := service.NewFriendRequestService()
	service.SetFriendRequestRepo(mockRepo)
	service.SetUserService(mockUserService)

	mockUserService.On("GetUserByID", "invalidUser").Return(nil, fmt.Errorf("user not found"))

	err := service.SendFriendRequest("invalidUser", "user2")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sender user not found")
}

func TestFriendRequestService_SendFriendRequest_InvalidRecipientID(t *testing.T) {
	mockRepo := new(tests.MockFriendRequestRepository)
	mockUserService := new(tests.MockUserService)
	service := service.NewFriendRequestService()
	service.SetFriendRequestRepo(mockRepo)
	service.SetUserService(mockUserService)

	mockUserService.On("GetUserByID", "user1").Return(&entity.User{ID: "user1"}, nil)
	mockUserService.On("GetUserByID", "invalidUser").Return(nil, fmt.Errorf("user not found"))

	err := service.SendFriendRequest("user1", "invalidUser")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "recipient user not found")
}

func TestFriendRequestService_SendFriendRequest_AlreadyExists(t *testing.T) {
	mockRepo := new(tests.MockFriendRequestRepository)
	mockUserService := new(tests.MockUserService)
	service := service.NewFriendRequestService()
	service.SetFriendRequestRepo(mockRepo)
	service.SetUserService(mockUserService)

	mockUserService.On("GetUserByID", "user1").Return(&entity.User{ID: "user1"}, nil)
	mockUserService.On("GetUserByID", "user2").Return(&entity.User{ID: "user2"}, nil)
	mockRepo.On("GetByUsers", "user1", "user2").Return(&tests.ValidFriendRequest, nil)

	err := service.SendFriendRequest("user1", "user2")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "friend request already exists")
}

func TestFriendRequestService_RespondToRequest_NonExistent(t *testing.T) {
	mockRepo := new(tests.MockFriendRequestRepository)
	service := service.NewFriendRequestService()
	service.SetFriendRequestRepo(mockRepo)

	mockRepo.On("GetByUsers", "user1", "user2").Return(nil, fmt.Errorf("not found"))

	err := service.RespondToFriendRequest("user1", "user2", true)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "friend request not found")
}

func TestFriendRequestService_RespondToRequest_AlreadyProcessed(t *testing.T) {
	mockRepo := new(tests.MockFriendRequestRepository)
	service := service.NewFriendRequestService()
	service.SetFriendRequestRepo(mockRepo)

	processedRequest := tests.AcceptedFriendRequest
	mockRepo.On("GetByUsers", "user1", "user2").Return(&processedRequest, nil)

	err := service.RespondToFriendRequest("user1", "user2", true)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "friend request already processed")
}

func TestFriendRequestService_SendFriendRequest_EmptyUserID(t *testing.T) {
	mockRepo := new(tests.MockFriendRequestRepository)
	mockUserService := new(tests.MockUserService)
	service := service.NewFriendRequestService()
	service.SetFriendRequestRepo(mockRepo)
	service.SetUserService(mockUserService)

	err := service.SendFriendRequest("", "user2")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user IDs cannot be empty")
}
