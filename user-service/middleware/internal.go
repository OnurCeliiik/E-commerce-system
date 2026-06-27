package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// RequireInternalSecret blocks requests without the shared service-to-service secret.
// Not exposed via the gateway — only reachable on the Docker network.
func RequireInternalSecret() gin.HandlerFunc {
	secret := os.Getenv("INTERNAL_SERVICE_SECRET")

	return func(c *gin.Context) {
		if secret == "" || c.GetHeader("X-Internal-Secret") != secret {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
