package main

import (
	"log"
	"os"
	"time"

	"github.com/OnurCeliiik/ecommerce/services/user/auth"
	"github.com/OnurCeliiik/ecommerce/services/user/database"
	"github.com/OnurCeliiik/ecommerce/services/user/handlers"
	"github.com/OnurCeliiik/ecommerce/services/user/repository"
	"github.com/OnurCeliiik/ecommerce/services/user/routes"
	"github.com/OnurCeliiik/ecommerce/services/user/service"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := database.MigrateDB(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	tokenProvider, err := auth.NewHMACProvider(os.Getenv("JWT_SECRET"), 24*time.Hour)
	if err != nil {
		log.Fatalf("failed to create token provider: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, tokenProvider)
	userHandler := handlers.NewUserHandler(userSvc)

	router := gin.Default()
	routes.RegisterRoutes(router, routes.Dependencies{
		DB:             db,
		UserHandler:    userHandler,
		TokenValidator: tokenProvider,
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
