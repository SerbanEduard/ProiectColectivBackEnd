package routes

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {
	userController := controller.NewUserController()

	r.POST("/users/signup", userController.SignUp)
	r.POST("/users/login", userController.Login)
	r.GET("/users/:id", userController.GetUser)
	r.GET("/users", userController.GetAllUsers)
	r.PATCH("/users/:id", controller.JWTAuthMiddleware(), controller.RequireOwner("id"), userController.UpdateUser)
	r.PUT("/users/:id/password", controller.JWTAuthMiddleware(), controller.RequireOwner("id"), userController.UpdateUserPassword)
	r.PUT("/users/:id/statistics", controller.JWTAuthMiddleware(), controller.RequireOwner("id"), userController.UpdateUserStatistics)
	r.DELETE("/users/:id", controller.JWTAuthMiddleware(), controller.RequireOwner("id"), userController.DeleteUser)
}
