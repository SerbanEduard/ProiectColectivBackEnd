package routes

import (
    "github.com/SerbanEduard/ProiectColectivBackEnd/controller"
    "github.com/gin-gonic/gin"
)

func SetupFriendRequestRoutes(r *gin.Engine) {
    friendRequestController := controller.NewFriendRequestController()

    // Protected endpoints - require JWT
    protected := r.Group("/")
    protected.Use(controller.JWTAuthMiddleware())
    {
        protected.POST("/friend-requests/:fromUserId/:toUserId", friendRequestController.SendFriendRequest)
        protected.PUT("/friend-requests/:fromUserId/:toUserId", friendRequestController.RespondToFriendRequest)
        protected.GET("/friend-requests/pending/:userId", friendRequestController.GetPendingRequests)
    }
}