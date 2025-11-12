package routes

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/gin-gonic/gin"
)

func VoiceRoutes(router *gin.Engine) {
	voiceController := controller.NewVoiceController()

	voice := router.Group("/voice")
	{
		voice.GET("/:teamId", voiceController.JoinVoiceRoom)
		voice.POST("/rooms/:teamId", voiceController.CreateVoiceRoom)
	}
}
