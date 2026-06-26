package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/OnurCeliiik/ecommerce/services/inventory/database"
	"github.com/OnurCeliiik/ecommerce/services/inventory/handlers"
	kafkasub "github.com/OnurCeliiik/ecommerce/services/inventory/kafka"
	"github.com/OnurCeliiik/ecommerce/services/inventory/repository"
	"github.com/OnurCeliiik/ecommerce/services/inventory/routes"
	"github.com/OnurCeliiik/ecommerce/services/inventory/service"
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

	inventoryRepo := repository.NewInventoryRepository(db)
	inventoryService := service.NewInventoryService(inventoryRepo)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	if brokers := strings.TrimSpace(os.Getenv("KAFKA_BROKERS")); brokers != "" {
		consumer, err := kafkasub.NewOrderEventConsumer(brokers, inventoryService)
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
		DB:               db,
		InventoryHandler: inventoryHandler,
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
