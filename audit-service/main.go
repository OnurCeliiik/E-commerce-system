package main

import (
	"context"
	"log"
	"time"

	"github.com/OnurCeliiik/ecommerce/services/audit/database"
	"github.com/OnurCeliiik/ecommerce/services/audit/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	mongo, err := database.ConnectFromEnv()
	if err != nil {
		log.Fatalf("failed to connect mongodb: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mongo.Close(ctx); err != nil {
			log.Printf("failed to close mongodb: %v", err)
		}
	}()

	log.Println("connected to MongoDB")

	router := gin.Default()
	routes.RegisterRoutes(router, routes.Dependencies{
		Mongo: mongo,
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start audit service: %v", err)
	}
}
