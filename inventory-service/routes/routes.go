package routes

import (
	"github.com/OnurCeliiik/ecommerce/services/inventory/handlers"
	"github.com/OnurCeliiik/ecommerce/services/inventory/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB               *gorm.DB
	InventoryHandler *handlers.InventoryHandler
}

func RegisterRoutes(router *gin.Engine, deps Dependencies) {
	router.Use(middleware.PrometheusMiddleware())
	router.GET("/health", handlers.HealthCheckHandler(deps.DB))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/inventory/:product_id", deps.InventoryHandler.GetInventory)
		v1.PUT("/inventory/:product_id", deps.InventoryHandler.UpdateInventory)
	}
}
