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

// TODO: Add Swagger comment
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

// TODO: Add Swagger comment
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

// TODO: Add swagger comment
func (fc *FriendRequestController) GetPendingRequests(c *gin.Context) {
	userID := c.Param("userId")

	requests, err := fc.friendRequestService.GetPendingRequests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}
