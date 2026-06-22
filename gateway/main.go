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

	tokenProvider, err := auth.NewHMACProvider(os.Getenv("JWT_SECRET"))
	if err != nil {
		log.Fatalf("failed to create token provider: %v", err)
	}

	userProxy, err := proxy.NewUserService(userServiceURL)
	if err != nil {
		log.Fatalf("failed to create user service proxy: %v", err)
	}

	router := gin.Default()
	routes.RegisterRoutes(router, routes.Dependencies{
		UserServiceProxy: userProxy,
		AuthMiddleware:   middleware.Auth(tokenProvider),
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to run gateway: %v", err)
	}
}
