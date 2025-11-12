package routes

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/gin-gonic/gin"
)

func SetupTeamRoutes(r *gin.Engine) {
	teamController := controller.NewTeamController()

	// Protected endpoints - require JWT
	protected := r.Group("/")
	protected.Use(controller.JWTAuthMiddleware())
	{
		protected.POST("/teams", teamController.NewTeam)
		protected.GET("/teams/:id", teamController.GetTeam)
		protected.GET("/teams", teamController.GetAllTeams)
		protected.PUT("/teams/:id", teamController.UpdateTeam)
		protected.DELETE("/teams/:id", teamController.DeleteTeam)
		protected.GET("/teams/search", teamController.GetXTeamsByPrefix)
		protected.GET("/teams/by-name", teamController.GetTeamsByName)
		protected.POST("/teams/addUserToTeam", teamController.AddUserToTeam)
	}
}
