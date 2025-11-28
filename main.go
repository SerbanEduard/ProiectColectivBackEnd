package main

import (
	"log"
	"os"

	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/SerbanEduard/ProiectColectivBackEnd/docs"
	"github.com/SerbanEduard/ProiectColectivBackEnd/routes"
	"github.com/gin-gonic/gin"
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
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	log.Println("Starting StudyWithMe API server...")
	log.Printf("Gin mode: %s", gin.Mode())

	config.InitFirebase()

	r := routes.SetupRoutes()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server starting on port 8080...")
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
