package persistence

import (
    "context"
    "fmt"
    "github.com/SerbanEduard/ProiectColectivBackEnd/config"
    "github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

type FriendRequestRepository struct{}

func NewFriendRequestRepository() *FriendRequestRepository {
    return &FriendRequestRepository{}
}

func (fr *FriendRequestRepository) Create(request *entity.FriendRequest) error {
    ctx := context.Background()
    ref := config.FirebaseDB.NewRef("friendRequests/" + request.Key())
    return ref.Set(ctx, request)
}

func (fr *FriendRequestRepository) GetByUsers(fromUserID, toUserID string) (*entity.FriendRequest, error) {
    ctx := context.Background()
    key := fromUserID + ":" + toUserID
    ref := config.FirebaseDB.NewRef("friendRequests/" + key)
    
    var request entity.FriendRequest
    if err := ref.Get(ctx, &request); err != nil {
        if err.Error() == "firebase: no such child" {
            return nil, fmt.Errorf("friend request not found")
        }
        return nil, fmt.Errorf("error retrieving friend request: %v", err)
    }
    
    return &request, nil
}

func (fr *FriendRequestRepository) Update(request *entity.FriendRequest) error {
    ctx := context.Background()
    ref := config.FirebaseDB.NewRef("friendRequests/" + request.Key())
    return ref.Set(ctx, request)
}

func (fr *FriendRequestRepository) GetPendingRequestsForUser(userID string) ([]*entity.FriendRequest, error) {
    ctx := context.Background()
    ref := config.FirebaseDB.NewRef("friendRequests")
    
    var requestsMap map[string]*entity.FriendRequest
    if err := ref.Get(ctx, &requestsMap); err != nil {
        return nil, err
    }
    
    var pendingRequests []*entity.FriendRequest
    for _, request := range requestsMap {
        if request.ToUserID == userID && request.Status == entity.PENDING {
            pendingRequests = append(pendingRequests, request)
        }
    }
    return pendingRequests, nil
}