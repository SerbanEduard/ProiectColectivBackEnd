package dto

import "time"

type RespondFriendRequestRequest struct {
	Accept bool `json:"accept"`
}

type FriendRequestResponse struct {
	FromUserID string    `json:"fromUserId"`
	ToUserID   string    `json:"toUserId"`
	CreatedAt  time.Time `json:"createdAt"`
	Status     string    `json:"status"`
}

type FriendRequestListResponse struct {
	Requests []FriendRequestResponse `json:"requests"`
}
