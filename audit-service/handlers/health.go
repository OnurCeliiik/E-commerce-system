package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/OnurCeliiik/ecommerce/services/audit/database"
	"github.com/gin-gonic/gin"
)

func HealthCheckHandler(mongo *database.MongoDBClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := mongo.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "degraded",
				"service": "audit",
				"error":   "mongodb unreachable",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "audit"})
	}
}
