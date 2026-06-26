package routes

import (
	"github.com/OnurCeliiik/ecommerce/services/notification/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB *gorm.DB
}

func RegisterRoutes(router *gin.Engine, deps Dependencies) {
	router.GET("/health", handlers.HealthCheckHandler(deps.DB))
}
