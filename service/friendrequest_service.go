package service

import (
	"fmt"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/persistence"
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

	GetFriendsForUser(userID string) ([]string, error)
}

type UserServiceInterface interface {
	GetUserByID(userID string) (*entity.User, error)
}

type FriendRequestServiceInterface interface {
	SendFriendRequest(fromUserID, toUserID string) error
	RespondToFriendRequest(fromUserID, toUserID string, accept bool) error
	GetPendingRequests(userID string) ([]*entity.FriendRequest, error)

	GetFriends(userID string) ([]*entity.User, error)

	GetMutualFriends(userA, userB string) ([]*entity.User, error)
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

	var (
		sender    *entity.User
		recipient *entity.User
		existing  *entity.FriendRequest
		err       error
	)

	sender, err = fs.userService.GetUserByID(fromUserID)
	if err != nil || sender == nil {
		return fmt.Errorf("sender user not found")
	}

	recipient, err = fs.userService.GetUserByID(toUserID)
	if err != nil || recipient == nil {
		return fmt.Errorf("recipient user not found")
	}

	existing, _ = fs.friendRequestRepo.GetByUsers(fromUserID, toUserID)
	if existing != nil {
		return fmt.Errorf("friend request already exists")
	}

	request := &entity.FriendRequest{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Status:     entity.PENDING,
	}

	if err := fs.friendRequestRepo.Create(request); err != nil {
		return fmt.Errorf("create friend request: %w", err)
	}

	return nil
}

func (fs *FriendRequestService) RespondToFriendRequest(fromUserID, toUserID string, accept bool) error {
	if fromUserID == "" || toUserID == "" {
		return fmt.Errorf("user IDs cannot be empty")
	}

	var (
		request *entity.FriendRequest
		err     error
	)

	request, err = fs.friendRequestRepo.GetByUsers(fromUserID, toUserID)
	if err != nil || request == nil {
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

	if err := fs.friendRequestRepo.Update(request); err != nil {
		return fmt.Errorf("update friend request: %w", err)
	}

	return nil
}

func (fs *FriendRequestService) GetPendingRequests(userID string) ([]*entity.FriendRequest, error) {
	reqs, err := fs.friendRequestRepo.GetPendingRequestsForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("get pending requests: %w", err)
	}
	return reqs, nil
}

func (fs *FriendRequestService) GetFriends(userID string) ([]*entity.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("user id cannot be empty")
	}

	friendIDs, err := fs.friendRequestRepo.GetFriendsForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("get friends: %w", err)
	}

	var friends []*entity.User
	for _, fid := range friendIDs {
		u, err := fs.userService.GetUserByID(fid)
		if err != nil || u == nil {
			// skip missing users
			continue
		}
		friends = append(friends, u)
	}
	return friends, nil
}

func (fs *FriendRequestService) GetMutualFriends(userA, userB string) ([]*entity.User, error) {
	if userA == "" || userB == "" {
		return nil, fmt.Errorf("user ids cannot be empty")
	}

	idsA, err := fs.friendRequestRepo.GetFriendsForUser(userA)
	if err != nil {
		return nil, fmt.Errorf("get friends for user %s: %w", userA, err)
	}
	idsB, err := fs.friendRequestRepo.GetFriendsForUser(userB)
	if err != nil {
		return nil, fmt.Errorf("get friends for user %s: %w", userB, err)
	}

	setA := make(map[string]struct{}, len(idsA))
	for _, id := range idsA {
		setA[id] = struct{}{}
	}

	var mutualIDs []string
	for _, id := range idsB {
		if _, ok := setA[id]; ok {
			mutualIDs = append(mutualIDs, id)
		}
	}

	var mutualUsers []*entity.User
	for _, id := range mutualIDs {
		u, err := fs.userService.GetUserByID(id)
		if err != nil || u == nil {
			continue
		}
		mutualUsers = append(mutualUsers, u)
	}

	return mutualUsers, nil
}

func (fs *FriendRequestService) SetFriendRequestRepo(repo FriendRequestRepositoryInterface) {
	fs.friendRequestRepo = repo
}

func (fs *FriendRequestService) SetUserService(service UserServiceInterface) {
	fs.userService = service
}
