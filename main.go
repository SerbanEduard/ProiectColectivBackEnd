package main

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/SerbanEduard/ProiectColectivBackEnd/routes"
)

func main() {
	config.InitFirebase()

	r := routes.SetupRoutes()

	r.Run(":8080")
}
