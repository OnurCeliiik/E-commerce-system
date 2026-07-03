package routes

import (
	"github.com/OnurCeliiik/ecommerce/services/audit/database"
	"github.com/OnurCeliiik/ecommerce/services/audit/handlers"
	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Mongo *database.MongoDBClient
}

func RegisterRoutes(router *gin.Engine, deps Dependencies) {
	router.GET("/health", handlers.HealthCheckHandler(deps.Mongo))
}
