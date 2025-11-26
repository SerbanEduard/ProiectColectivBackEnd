package routes

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/gin-gonic/gin"
)

// VoiceRoutes sets up all the API routes for voice chat functionality.
func VoiceRoutes(router *gin.Engine) {
	voiceController := controller.NewVoiceController()

	voice := router.Group("/voice")
	voice.Use(controller.JWTAuthMiddleware())
	{
		voice.GET("/join/:roomId", voiceController.JoinVoiceRoom)
		voice.POST("/rooms/:teamId", voiceController.CreateVoiceRoom)
		voice.GET("/rooms/:teamId", voiceController.GetActiveRooms)
		voice.POST("/private/call", voiceController.StartPrivateCall)
		voice.GET("/joinable", voiceController.GetJoinableRooms)
	}
}
