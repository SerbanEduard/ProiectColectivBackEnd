package main

import (
	"time"

	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/SerbanEduard/ProiectColectivBackEnd/docs"
	"github.com/SerbanEduard/ProiectColectivBackEnd/routes"
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title						StudyWithMe API
// @version					1.0
// @BasePath					/
// @SecurityDefinitions.apiKey	Bearer
// @in							header
// @name						Authorization
// @description				Bearer JWT token. Example: "Bearer YOUR_TOKEN"
func main() {
	config.InitFirebase()

	r := routes.SetupRoutes()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://studyflow-6qwx.onrender.co"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup Swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
