package routes

import (
    "github.com/SerbanEduard/ProiectColectivBackEnd/controller"
    "github.com/gin-gonic/gin"
)

func SetupFriendRequestRoutes(r *gin.Engine) {
    friendRequestController := controller.NewFriendRequestController()

    r.POST("/friend-requests/:fromUserId/:toUserId", friendRequestController.SendFriendRequest)
    r.PUT("/friend-requests/:fromUserId/:toUserId", friendRequestController.RespondToFriendRequest)
    r.GET("/friend-requests/pending/:userId", friendRequestController.GetPendingRequests)
}