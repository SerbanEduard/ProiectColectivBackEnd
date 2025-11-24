package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://studyflow-6qwx.onrender.com", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from ProiectColectivBackEnd!",
		})
	})

	SetupUserRoutes(r)
	SetupTeamRoutes(r)
	FileRoutes(r)
	SetupMessageRoutes(r)
	SetupFriendRequestRoutes(r)
	VoiceRoutes(r)
	SetupQuizRoutes(r)

	return r
}
