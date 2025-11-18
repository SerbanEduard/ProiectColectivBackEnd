package routes

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/gin-gonic/gin"
)

func SetupQuizRoutes(r *gin.Engine) {
	quizController := controller.NewQuizController()

	protected := r.Group("/")
	protected.Use(controller.JWTAuthMiddleware())
	{
		protected.POST("/quizzes", quizController.CreateQuiz)
	}
}
