package routes

import (
	"github.com/OnurCeliiik/ecommerce/services/order/handlers"
	"github.com/OnurCeliiik/ecommerce/services/order/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB             *gorm.DB
	OrderHandler   *handlers.OrderHandler
	AuthMiddleware gin.HandlerFunc
}

func RegisterRoutes(router *gin.Engine, deps Dependencies) {
	router.Use(middleware.PrometheusMiddleware())
	router.GET("/health", handlers.HealthCheckHandler(deps.DB))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := router.Group("/api/v1")
	{
		protected := v1.Group("")
		protected.Use(deps.AuthMiddleware)
		protected.POST("/orders", deps.OrderHandler.CreateOrder)
		protected.GET("/orders/:id", deps.OrderHandler.GetOrder)
		protected.GET("/orders/me", deps.OrderHandler.GetOrders)
	}
}
