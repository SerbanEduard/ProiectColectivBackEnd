package routes

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/gin-gonic/gin"
)

func SetupTeamRequestRoutes(r *gin.Engine) {
	trController := controller.NewTeamRequestController()
	protected := r.Group("/")
	r.POST("/teamRequests", trController.CreateTeamRequest)
	r.PUT("/teamRequests/:id/accept", trController.AcceptTeamRequest)
	r.DELETE("/teamRequests/:id/reject", trController.RejectTeamRequest)
	r.GET("/teamRequests", trController.GetAllTeamRequests)
	r.GET("/teamRequests/user/:userId", trController.GetTeamRequestsByUser)
	protected.Use(controller.JWTAuthMiddleware())
	{
		/*protected.POST("/teamRequests", trController.CreateTeamRequest)
		protected.PUT("/teamRequests/:id/accept", trController.AcceptTeamRequest)
		protected.DELETE("/teamRequests/:id/reject", trController.RejectTeamRequest)
		protected.GET("/teamRequests", trController.GetAllTeamRequests)
		protected.GET("/teamRequests/user/:userId", trController.GetTeamRequestsByUser)*/
	}
}
