package controller

import (
	"net/http"

	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/gin-gonic/gin"
)

type FriendRequestController struct {
	friendRequestService service.FriendRequestServiceInterface
}

func NewFriendRequestController() *FriendRequestController {
	return &FriendRequestController{
		friendRequestService: service.NewFriendRequestService(),
	}
}

func (fc *FriendRequestController) SetFriendRequestService(service service.FriendRequestServiceInterface) {
	fc.friendRequestService = service
}

// @Summary		Send a friend request
// @Description	Send a friend request from one user to another
// @Tags			default
// @Param			fromUserId	path		string	true	"Sender User ID"
// @Param			toUserId	path		string	true	"Recipient User ID"
// @Success		201			{object}	nil
// @Failure		400			{object}	map[string]string
// @Failure		500			{object}	map[string]string
// @Router			/friend-requests/{fromUserId}/{toUserId} [post]
func (fc *FriendRequestController) SendFriendRequest(c *gin.Context) {
	fromUserID := c.Param("fromUserId")
	toUserID := c.Param("toUserId")

	if fromUserID == "" || toUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user IDs"})
		return
	}

	err := fc.friendRequestService.SendFriendRequest(fromUserID, toUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Friend request sent successfully"})
}

// @Summary		Respond to a friend request
// @Description	Accept or deny a friend request
// @Tags			default
// @Param			fromUserId	path		string							true	"Sender User ID"
// @Param			toUserId	path		string							true	"Recipient User ID"
// @Param			body		body		dto.RespondFriendRequestRequest	true	"Accept or deny"
// @Success		200			{object}	nil
// @Failure		400			{object}	map[string]string
// @Failure		404			{object}	map[string]string
// @Failure		500			{object}	map[string]string
// @Router			/friend-requests/{fromUserId}/{toUserId} [put]
func (fc *FriendRequestController) RespondToFriendRequest(c *gin.Context) {
	fromUserID := c.Param("fromUserId")
	toUserID := c.Param("toUserId")

	if fromUserID == "" || toUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user IDs"})
		return
	}

	var request struct {
		Accept bool `json:"accept"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := fc.friendRequestService.RespondToFriendRequest(fromUserID, toUserID, request.Accept)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request processed successfully"})
}

// @Summary		Get pending friend requests
// @Description	Get pending friend requests for a user
// @Tags			default
// @Param			userId	path		string	true	"User ID"
// @Success		200		{object}	dto.FriendRequestListResponse
// @Failure		500		{object}	map[string]string
// @Router			/friend-requests/{userId} [get]
func (fc *FriendRequestController) GetPendingRequests(c *gin.Context) {
	userID := c.Param("userId")

	requests, err := fc.friendRequestService.GetPendingRequests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}
