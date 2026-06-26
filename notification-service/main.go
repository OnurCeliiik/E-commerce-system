package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/OnurCeliiik/ecommerce/services/notification/database"
	emailsender "github.com/OnurCeliiik/ecommerce/services/notification/email"
	kafkasub "github.com/OnurCeliiik/ecommerce/services/notification/kafka"
	"github.com/OnurCeliiik/ecommerce/services/notification/routes"
	"github.com/OnurCeliiik/ecommerce/services/notification/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := database.MigrateDB(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	emailSender := emailsender.NewLogSender()
	notificationService := service.NewNotificationService(emailSender)

	if brokers := strings.TrimSpace(os.Getenv("KAFKA_BROKERS")); brokers != "" {
		consumer, err := kafkasub.NewOrderEventConsumer(brokers, notificationService)
		if err != nil {
			log.Fatalf("failed to create kafka consumer: %v", err)
		}
		go func() {
			if err := consumer.Run(context.Background()); err != nil {
				log.Printf("kafka consumer stopped: %v", err)
			}
		}()
	}

	router := gin.Default()
	routes.RegisterRoutes(router, routes.Dependencies{
		DB: db,
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
