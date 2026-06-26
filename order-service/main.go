package main

import (
	"log"
	"os"

	"github.com/OnurCeliiik/ecommerce/services/order/auth"
	"github.com/OnurCeliiik/ecommerce/services/order/catalog"
	"github.com/OnurCeliiik/ecommerce/services/order/database"
	"github.com/OnurCeliiik/ecommerce/services/order/handlers"
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
	orderSvc := service.NewOrderService(orderRepo, catalogClient)
	orderHandler := handlers.NewOrderHandler(orderSvc)

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
