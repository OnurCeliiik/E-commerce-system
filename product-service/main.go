package main

import (
	"log"

	"github.com/OnurCeliiik/ecommerce/services/product/database"
	"github.com/OnurCeliiik/ecommerce/services/product/handlers"
	"github.com/OnurCeliiik/ecommerce/services/product/repository"
	"github.com/OnurCeliiik/ecommerce/services/product/routes"
	"github.com/OnurCeliiik/ecommerce/services/product/service"
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

	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productSvc)

	router := gin.Default()
	routes.RegisterRoutes(router, routes.Dependencies{
		DB:             db,
		ProductHandler: productHandler,
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
