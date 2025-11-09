package service

import (
	"fmt"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/persistence"
	"github.com/SerbanEduard/ProiectColectivBackEnd/validator"
)

type FriendRequestService struct {
	friendRequestRepo FriendRequestRepositoryInterface
	userService       UserServiceInterface
}

type FriendRequestRepositoryInterface interface {
	Create(request *entity.FriendRequest) error
	GetByUsers(fromUserID, toUserID string) (*entity.FriendRequest, error)
	Update(request *entity.FriendRequest) error
	GetPendingRequestsForUser(userID string) ([]*entity.FriendRequest, error)
}

type UserServiceInterface interface {
	GetUserByID(userID string) (*entity.User, error)
}

type FriendRequestServiceInterface interface {
	SendFriendRequest(fromUserID, toUserID string) error
	RespondToFriendRequest(fromUserID, toUserID string, accept bool) error
	GetPendingRequests(userID string) ([]*entity.FriendRequest, error)
}

func NewFriendRequestService() *FriendRequestService {
	return &FriendRequestService{
		friendRequestRepo: persistence.NewFriendRequestRepository(),
		userService:       NewUserService(),
	}
}

func (fs *FriendRequestService) SendFriendRequest(fromUserID, toUserID string) error {
	if fromUserID == "" || toUserID == "" {
		return fmt.Errorf("user IDs cannot be empty")
	}

	if err := validator.ValidateFriendRequest(fromUserID, toUserID); err != nil {
		return err
	}

	if _, err := fs.userService.GetUserByID(fromUserID); err != nil {
		return fmt.Errorf("sender user not found")
	}
	if _, err := fs.userService.GetUserByID(toUserID); err != nil {
		return fmt.Errorf("recipient user not found")
	}

	if _, err := fs.friendRequestRepo.GetByUsers(fromUserID, toUserID); err == nil {
		return fmt.Errorf("friend request already exists")
	}

	request := entity.NewFriendRequest(fromUserID, toUserID)
	return fs.friendRequestRepo.Create(request)
}

func (fs *FriendRequestService) RespondToFriendRequest(fromUserID, toUserID string, accept bool) error {
	request, err := fs.friendRequestRepo.GetByUsers(fromUserID, toUserID)
	if err != nil {
		return fmt.Errorf("friend request not found")
	}

	if request.Status != entity.PENDING {
		return fmt.Errorf("friend request already processed")
	}

	if accept {
		request.Status = entity.ACCEPTED
	} else {
		request.Status = entity.DENIED
	}

	return fs.friendRequestRepo.Update(request)
}

func (fs *FriendRequestService) GetPendingRequests(userID string) ([]*entity.FriendRequest, error) {
	return fs.friendRequestRepo.GetPendingRequestsForUser(userID)
}

// Add these methods for testing
func (fs *FriendRequestService) SetFriendRequestRepo(repo FriendRequestRepositoryInterface) {
	fs.friendRequestRepo = repo
}

func (fs *FriendRequestService) SetUserService(service UserServiceInterface) {
	fs.userService = service
}
