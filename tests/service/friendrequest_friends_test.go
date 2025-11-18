package service_test

import (
	"testing"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/stretchr/testify/assert"
)

func TestFriendRequestService_GetFriends_Success(t *testing.T) {
	mockRepo := new(tests.MockFriendRequestRepository)
	mockUserService := new(tests.MockUserService)

	svc := service.NewFriendRequestService()
	svc.SetFriendRequestRepo(mockRepo)
	svc.SetUserService(mockUserService)

	mockRepo.On("GetFriendsForUser", "user1").Return([]string{"f1", "f2"}, nil)
	mockUserService.On("GetUserByID", "f1").Return(&entity.User{ID: "f1", Username: "u1"}, nil)
	mockUserService.On("GetUserByID", "f2").Return(&entity.User{ID: "f2", Username: "u2"}, nil)

	friends, err := svc.GetFriends("user1")

	assert.NoError(t, err)
	assert.Len(t, friends, 2)
	mockRepo.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

func TestFriendRequestService_GetMutualFriends_Success(t *testing.T) {
	mockRepo := new(tests.MockFriendRequestRepository)
	mockUserService := new(tests.MockUserService)

	svc := service.NewFriendRequestService()
	svc.SetFriendRequestRepo(mockRepo)
	svc.SetUserService(mockUserService)

	mockRepo.On("GetFriendsForUser", "a").Return([]string{"f1", "f2", "f3"}, nil)
	mockRepo.On("GetFriendsForUser", "b").Return([]string{"f2", "f3", "f4"}, nil)
	mockUserService.On("GetUserByID", "f2").Return(&entity.User{ID: "f2", Username: "u2"}, nil)
	mockUserService.On("GetUserByID", "f3").Return(&entity.User{ID: "f3", Username: "u3"}, nil)

	mutual, err := svc.GetMutualFriends("a", "b")

	assert.NoError(t, err)
	assert.Len(t, mutual, 2)
	mockRepo.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}
