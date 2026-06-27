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
		reservedConsumer, err := kafkasub.NewInventoryReservedConsumer(brokers, notificationService)
		if err != nil {
			log.Fatalf("failed to create inventory.reserved consumer: %v", err)
		}
		failedConsumer, err := kafkasub.NewInventoryReservationFailedConsumer(brokers, notificationService)
		if err != nil {
			log.Fatalf("failed to create inventory.reservation_failed consumer: %v", err)
		}
		go func() {
			if err := reservedConsumer.Run(context.Background()); err != nil {
				log.Printf("inventory.reserved consumer stopped: %v", err)
			}
		}()
		go func() {
			if err := failedConsumer.Run(context.Background()); err != nil {
				log.Printf("inventory.reservation_failed consumer stopped: %v", err)
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
