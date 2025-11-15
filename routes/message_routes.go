package routes

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/gin-gonic/gin"
)

func SetupMessageRoutes(r *gin.Engine) {
	messageController := controller.NewMessageController()

	// Protected endpoints
	protected := r.Group("/")
	protected.Use(controller.JWTAuthMiddleware())
	{
		protected.GET("/messages/connect", messageController.Connect)
		protected.POST("/messages/direct", messageController.NewDirectMessage)
		protected.POST("/messages/team", messageController.NewTeamMessage)

		// TODO: get by id, get direct/team messages, edit messages
	}
}
