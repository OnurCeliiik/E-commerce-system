package main

import (
	"log"
	"os"

	"github.com/OnurCeliiik/ecommerce/gateway/auth"
	"github.com/OnurCeliiik/ecommerce/gateway/middleware"
	"github.com/OnurCeliiik/ecommerce/gateway/proxy"
	"github.com/OnurCeliiik/ecommerce/gateway/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		log.Fatal("USER_SERVICE_URL is not set")
	}

	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		log.Fatal("PRODUCT_SERVICE_URL is not set")
	}

	inventoryServiceURL := os.Getenv("INVENTORY_SERVICE_URL")
	if inventoryServiceURL == "" {
		log.Fatal("INVENTORY_SERVICE_URL is not set")
	}

	orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
	if orderServiceURL == "" {
		log.Fatal("ORDER_SERVICE_URL is not set")
	}

	tokenProvider, err := auth.NewHMACProvider(os.Getenv("JWT_SECRET"))
	if err != nil {
		log.Fatalf("failed to create token provider: %v", err)
	}

	userProxy, err := proxy.NewUserService(userServiceURL)
	if err != nil {
		log.Fatalf("failed to create user service proxy: %v", err)
	}

	productProxy, err := proxy.NewProductService(productServiceURL)
	if err != nil {
		log.Fatalf("failed to create product service proxy: %v", err)
	}

	inventoryProxy, err := proxy.NewInventoryService(inventoryServiceURL)
	if err != nil {
		log.Fatalf("failed to create inventory service proxy: %v", err)
	}

	orderProxy, err := proxy.NewOrderService(orderServiceURL)
	if err != nil {
		log.Fatalf("failed to create order service proxy: %v", err)
	}

	router := gin.Default()
	routes.RegisterRoutes(router, routes.Dependencies{
		UserServiceProxy:      userProxy,
		ProductServiceProxy:   productProxy,
		InventoryServiceProxy: inventoryProxy,
		OrderServiceProxy:     orderProxy,
		AuthMiddleware:        middleware.Auth(tokenProvider),
		RequireAdmin:          middleware.RequireRole(middleware.RoleAdmin),
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to run gateway: %v", err)
	}
}
