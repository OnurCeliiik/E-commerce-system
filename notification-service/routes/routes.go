package routes

import (
	"github.com/OnurCeliiik/ecommerce/services/notification/handlers"
	"github.com/OnurCeliiik/ecommerce/services/notification/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB *gorm.DB
}

func RegisterRoutes(router *gin.Engine, deps Dependencies) {

	router.Use(middleware.PrometheusMiddleware())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/health", handlers.HealthCheckHandler(deps.DB))
}
