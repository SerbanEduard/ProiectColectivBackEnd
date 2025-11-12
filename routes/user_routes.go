package routes

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {
	userController := controller.NewUserController()

	// Public endpoints
	r.POST("/users/signup", userController.SignUp)
	r.POST("/users/login", userController.Login)

	// Protected endpoints
	protected := r.Group("/")
	protected.Use(controller.JWTAuthMiddleware())
	{
		protected.GET("/users/:id", userController.GetUser)
		protected.GET("/users", userController.GetAllUsers)
		protected.PUT("/users/:id", userController.UpdateUser)
		protected.PUT("/users/:id/statistics", userController.UpdateUserStatistics)
		protected.DELETE("/users/:id", userController.DeleteUser)
	}
}
