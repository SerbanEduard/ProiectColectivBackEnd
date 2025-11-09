package routes

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {
	userController := controller.NewUserController()

	r.POST("/users/signup", userController.SignUp)
	r.GET("/users/:id", userController.GetUser)
	r.GET("/users", userController.GetAllUsers)
	r.PUT("/users/:id", userController.UpdateUser)
	r.PUT("/users/:id/statistics", userController.UpdateUserStatistics)
	r.DELETE("/users/:id", userController.DeleteUser)
}