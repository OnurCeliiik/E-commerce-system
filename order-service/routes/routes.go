package routes

import (
	"github.com/OnurCeliiik/ecommerce/services/order/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB              *gorm.DB
	OrderHandler    *handlers.OrderHandler
	AuthMiddleware  gin.HandlerFunc
}

func RegisterRoutes(router *gin.Engine, deps Dependencies) {
	router.GET("/health", handlers.HealthCheckHandler(deps.DB))

	v1 := router.Group("/api/v1")
	{
		protected := v1.Group("")
		protected.Use(deps.AuthMiddleware)
		protected.POST("/orders", deps.OrderHandler.CreateOrder)
	}
}
