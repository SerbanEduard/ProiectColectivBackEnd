package routes

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/gin-gonic/gin"
)

func SetupTeamRoutes(r *gin.Engine) {
	teamController := controller.NewTeamController()

	r.PUT("/teams/addUserToTeam", teamController.AddUserToTeam)              // Add a user to a team
	r.DELETE("/teams/deleteUserFromTeam", teamController.DeleteUserFromTeam) // Delete a user from a team

	r.GET("/teams/search", teamController.GetXTeamsByPrefix) // Query ?prefix=&limit=
	r.GET("/teams/by-name", teamController.GetTeamsByName)   // Query ?name=

	r.POST("/teams", teamController.NewTeam)          // Create a team
	r.GET("/teams/:id", teamController.GetTeam)       // Get a team by ID
	r.GET("/teams", teamController.GetAllTeams)       // Get all teams
	r.PUT("/teams/:id", teamController.UpdateTeam)    // Update a team
	r.DELETE("/teams/:id", teamController.DeleteTeam) // Delete a team

}
