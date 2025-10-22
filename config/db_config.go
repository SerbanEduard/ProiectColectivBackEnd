package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var FirebaseDB *db.Client

func InitFirebase() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	ctx := context.Background()

	credentialsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	databaseURL := os.Getenv("FIREBASE_DATABASE_URL")

	opt := option.WithCredentialsFile(credentialsPath)
	config := &firebase.Config{
		DatabaseURL: databaseURL,
	}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}

	FirebaseDB, err = app.Database(ctx)
	if err != nil {
		log.Fatalf("Error getting database client: %v", err)
	}

	log.Println("Firebase initialized successfully")
}
