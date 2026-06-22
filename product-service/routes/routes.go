package routes

import (
	"github.com/OnurCeliiik/ecommerce/services/product/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB              *gorm.DB
	ProductHandler  *handlers.ProductHandler
}

func RegisterRoutes(router *gin.Engine, deps Dependencies) {
	router.GET("/health", handlers.HealthCheck(deps.DB))

	v1 := router.Group("/api/v1")
	{
		v1.POST("/products", deps.ProductHandler.CreateProduct)
		v1.GET("/products", deps.ProductHandler.ListProducts)
		v1.GET("/products/:id", deps.ProductHandler.GetProductByID)
	}
}
