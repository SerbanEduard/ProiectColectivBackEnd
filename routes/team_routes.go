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
		protected.PUT("/teams/addUserToTeam", teamController.AddUserToTeam)              // Add a user to a team
		protected.DELETE("/teams/deleteUserFromTeam", teamController.DeleteUserFromTeam) // Delete a user from a team

		protected.GET("/teams/search", teamController.GetXTeamsByPrefix) // Query ?prefix=&limit=
		protected.GET("/teams/by-name", teamController.GetTeamsByName)   // Query ?name=

		protected.POST("/teams", teamController.NewTeam)          // Create a team
		protected.GET("/teams/:id", teamController.GetTeam)       // Get a team by ID
		protected.GET("/teams", teamController.GetAllTeams)       // Get all teams
		protected.PUT("/teams/:id", teamController.UpdateTeam)    // Update a team
		protected.DELETE("/teams/:id", teamController.DeleteTeam) // Delete a team
	}
}
