package routes

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/gin-gonic/gin"
)

func VoiceRoutes(router *gin.Engine) {
	voiceController := controller.NewVoiceController()

	// Protected endpoints - require JWT
	voice := router.Group("/voice")
	voice.Use(controller.JWTAuthMiddleware())
	{
		voice.GET("/:teamId", voiceController.JoinVoiceRoom)
		voice.POST("/rooms/:teamId", voiceController.CreateVoiceRoom)
		voice.DELETE("/:teamId/leave", voiceController.LeaveVoiceRoom)
	}
}
