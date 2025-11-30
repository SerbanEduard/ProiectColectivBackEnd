package routes

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/gin-gonic/gin"
)

func SetupEventRoutes(r *gin.Engine) {
	eventController := controller.NewEventController()

	// Protected endpoints
	protected := r.Group("/")
	protected.Use(controller.JWTAuthMiddleware())
	{
		protected.POST("/events", eventController.NewEvent)
		protected.GET("/events", eventController.GetEvents)
		protected.GET("/events/:id", eventController.GetEvent)
		protected.DELETE("/events/:id", eventController.DeleteEvent)
		protected.PATCH("/events/:id", eventController.UpdateEventDetails)
		protected.PATCH("/events/:id/status", eventController.UpdateUserStatus)
	}
}
