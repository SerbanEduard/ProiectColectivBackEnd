package routes

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/gin-gonic/gin"
)

func FileRoutes(router *gin.Engine) {
	fileController := controller.NewFileController()

	router.GET("/files", fileController.GetAllFiles)
	router.GET("/files/:id", fileController.GetFile)

	files := router.Group("/files")
	files.Use(controller.JWTAuthMiddleware())
	{
		files.POST("", fileController.UploadFile)
		files.DELETE(":id", fileController.DeleteFile)
	}
}
