package entity

import "time"

type FriendRequestStatus string

const (
	PENDING  FriendRequestStatus = "PENDING"
	ACCEPTED FriendRequestStatus = "ACCEPTED"
	DENIED   FriendRequestStatus = "DENIED"
)

type FriendRequest struct {
	FromUserID string              `json:"fromUserId"`
	ToUserID   string              `json:"toUserId"`
	CreatedAt  time.Time           `json:"createdAt"`
	Status     FriendRequestStatus `json:"status"`
}

func NewFriendRequest(fromUserID, toUserID string) *FriendRequest {
	return &FriendRequest{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		CreatedAt:  time.Now().UTC(),
		Status:     PENDING,
	}
}

func (fr *FriendRequest) Key() string {
	return fr.FromUserID + ":" + fr.ToUserID
}
