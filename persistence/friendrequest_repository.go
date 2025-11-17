package persistence

import (
	"context"
	"errors"
	"fmt"

	"firebase.google.com/go/v4/errorutils"

	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

const friendRequestsPath = "friendRequests"

var ErrFriendRequestNotFound = errors.New("friend request not found")

type FriendRequestRepository struct{}

func NewFriendRequestRepository() *FriendRequestRepository {
	return &FriendRequestRepository{}
}

func (fr *FriendRequestRepository) Create(request *entity.FriendRequest) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(friendRequestsPath + "/" + request.Key())
	if err := ref.Set(ctx, request); err != nil {
		return fmt.Errorf("create friend request: %w", err)
	}
	return nil
}

func (fr *FriendRequestRepository) GetByUsers(fromUserID, toUserID string) (*entity.FriendRequest, error) {
	ctx := context.Background()
	key := fromUserID + ":" + toUserID
	ref := config.FirebaseDB.NewRef(friendRequestsPath + "/" + key)

	var request entity.FriendRequest
	if err := ref.Get(ctx, &request); err != nil {
		if errorutils.IsNotFound(err) {
			return nil, fmt.Errorf("%w", ErrFriendRequestNotFound)
		}

		return nil, fmt.Errorf("get friend request %s: %w", key, err)
	}

	return &request, nil
}

func (fr *FriendRequestRepository) Update(request *entity.FriendRequest) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(friendRequestsPath + "/" + request.Key())
	if err := ref.Set(ctx, request); err != nil {
		return fmt.Errorf("update friend request: %w", err)
	}
	return nil
}

func (fr *FriendRequestRepository) GetPendingRequestsForUser(userID string) ([]*entity.FriendRequest, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(friendRequestsPath)

	var requestsMap map[string]*entity.FriendRequest
	if err := ref.Get(ctx, &requestsMap); err != nil {
		if errorutils.IsNotFound(err) {
			return []*entity.FriendRequest{}, nil
		}
		return nil, fmt.Errorf("get pending friend requests for user %s: %w", userID, err)
	}

	var pendingRequests []*entity.FriendRequest
	for _, request := range requestsMap {
		if request != nil && request.ToUserID == userID && request.Status == entity.PENDING {
			pendingRequests = append(pendingRequests, request)
		}
	}

	return pendingRequests, nil
}

func (fr *FriendRequestRepository) GetFriendsForUser(userID string) ([]string, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(friendRequestsPath)

	var requestsMap map[string]*entity.FriendRequest
	if err := ref.Get(ctx, &requestsMap); err != nil {
		if errorutils.IsNotFound(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("get friend requests map: %w", err)
	}

	friendSet := make(map[string]struct{})
	for _, req := range requestsMap {
		if req == nil {
			continue
		}
		if req.Status != entity.ACCEPTED {
			continue
		}
		if req.FromUserID == userID {
			friendSet[req.ToUserID] = struct{}{}
		} else if req.ToUserID == userID {
			friendSet[req.FromUserID] = struct{}{}
		}
	}

	friends := make([]string, 0, len(friendSet))
	for id := range friendSet {
		friends = append(friends, id)
	}
	return friends, nil
}
