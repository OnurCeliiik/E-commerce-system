package routes

import (
	"github.com/OnurCeliiik/ecommerce/services/product/handlers"
	"github.com/OnurCeliiik/ecommerce/services/product/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB             *gorm.DB
	ProductHandler *handlers.ProductHandler
}

func RegisterRoutes(router *gin.Engine, deps Dependencies) {
	router.Use(middleware.PrometheusMiddleware())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/health", handlers.HealthCheck(deps.DB))

	v1 := router.Group("/api/v1")
	{
		v1.POST("/products", deps.ProductHandler.CreateProduct)
		v1.GET("/products", deps.ProductHandler.ListProducts)
		v1.GET("/products/:id", deps.ProductHandler.GetProductByID)
		v1.PUT("/products/:id", deps.ProductHandler.UpdateProduct)
		v1.DELETE("/products/:id", deps.ProductHandler.DeleteProduct)
	}
}
