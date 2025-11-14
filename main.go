package main

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/SerbanEduard/ProiectColectivBackEnd/docs"
	"github.com/SerbanEduard/ProiectColectivBackEnd/routes"
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

	// Setup Swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
