package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/OnurCeliiik/ecommerce/services/order/auth"
	"github.com/OnurCeliiik/ecommerce/services/order/catalog"
	"github.com/OnurCeliiik/ecommerce/services/order/database"
	"github.com/OnurCeliiik/ecommerce/services/order/handlers"
	kafkapub "github.com/OnurCeliiik/ecommerce/services/order/kafka"
	"github.com/OnurCeliiik/ecommerce/services/order/middleware"
	"github.com/OnurCeliiik/ecommerce/services/order/repository"
	"github.com/OnurCeliiik/ecommerce/services/order/routes"
	"github.com/OnurCeliiik/ecommerce/services/order/service"
	"github.com/gin-gonic/gin"
)

func main() {

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := database.MigrateDB(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	serviceURL := os.Getenv("PRODUCT_SERVICE_URL")
	secret := os.Getenv("JWT_SECRET")

	orderRepo := repository.NewOrderRepository(db)
	catalogClient := catalog.NewHTTPProductClient(serviceURL)

	var publisher service.OrderEventPublisher = kafkapub.NoopPublisher{}
	if brokers := strings.TrimSpace(os.Getenv("KAFKA_BROKERS")); brokers != "" {
		kafkaPublisher, err := kafkapub.NewOrderEventPublisher(brokers)
		if err != nil {
			log.Fatalf("failed to create kafka publisher: %v", err)
		}
		publisher = kafkaPublisher
	}

	orderSvc := service.NewOrderService(orderRepo, catalogClient, publisher)
	orderHandler := handlers.NewOrderHandler(orderSvc)

	if brokers := strings.TrimSpace(os.Getenv("KAFKA_BROKERS")); brokers != "" {
		reservedConsumer, err := kafkapub.NewInventoryReservedConsumer(brokers, orderSvc)
		if err != nil {
			log.Fatalf("failed to create inventory.reserved consumer: %v", err)
		}
		failedConsumer, err := kafkapub.NewInventoryReservationFailedConsumer(brokers, orderSvc)
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

	tokenProvider, err := auth.NewHMACProvider(secret)
	if err != nil {
		log.Fatalf("failed to create token provider: %v", err)
	}

	router := gin.Default()
	routes.RegisterRoutes(router, routes.Dependencies{
		DB:             db,
		OrderHandler:   orderHandler,
		AuthMiddleware: middleware.Auth(tokenProvider),
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
