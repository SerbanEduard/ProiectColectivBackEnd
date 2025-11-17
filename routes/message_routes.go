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
		protected.POST("/messages", messageController.NewMessage)
		protected.GET("/messages", messageController.GetMessages)
		protected.GET("/messages/:id", messageController.GetMessage)
		protected.GET("/messages/connect", messageController.Connect)

		// TODO: edit messages
	}
}
