package routes

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/gin-gonic/gin"
)

func FileRoutes(router *gin.Engine) {
	fileController := controller.NewFileController()

	// All file endpoints are under /teams/:id/files (requires JWT)
	teams := router.Group("/teams")
	teams.Use(controller.JWTAuthMiddleware())
	{
		teams.GET("/:id/files", fileController.GetFilesByTeam)
		teams.POST("/:id/files", fileController.UploadFile)
		teams.GET("/:id/files/:fileId", fileController.GetFile)
		teams.DELETE("/:id/files/:fileId", fileController.DeleteFile)
	}
}
