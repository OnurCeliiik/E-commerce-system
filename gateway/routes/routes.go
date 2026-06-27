package routes

import (
	"net/http"

	"github.com/OnurCeliiik/ecommerce/gateway/proxy"
	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	UserServiceProxy      *proxy.UserService
	ProductServiceProxy   *proxy.ProductService
	InventoryServiceProxy *proxy.InventoryService
	OrderServiceProxy     *proxy.OrderService
	AuthMiddleware        gin.HandlerFunc
	RequireAdmin          gin.HandlerFunc
}

func RegisterRoutes(router *gin.Engine, deps Dependencies) {
	router.GET("/health", healthCheck)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", deps.UserServiceProxy.ServeHTTP)
		v1.POST("/login", deps.UserServiceProxy.ServeHTTP)

		v1.GET("/products", deps.ProductServiceProxy.ServeHTTP)
		v1.GET("/products/:id", deps.ProductServiceProxy.ServeHTTP)
		v1.GET("/inventory/:product_id", deps.InventoryServiceProxy.ServeHTTP)

		protected := v1.Group("")
		protected.Use(deps.AuthMiddleware)
		protected.GET("/me", deps.UserServiceProxy.ServeHTTP)
		protected.POST("/orders", deps.OrderServiceProxy.ServeHTTP)
		protected.GET("/orders/:id", deps.OrderServiceProxy.ServeHTTP)
		protected.GET("/orders/me", deps.OrderServiceProxy.ServeHTTP)

		admin := v1.Group("")
		admin.Use(deps.AuthMiddleware, deps.RequireAdmin)
		admin.POST("/products", deps.ProductServiceProxy.ServeHTTP)
		admin.PUT("/products/:id", deps.ProductServiceProxy.ServeHTTP)
		admin.DELETE("/products/:id", deps.ProductServiceProxy.ServeHTTP)
		admin.PUT("/inventory/:product_id", deps.InventoryServiceProxy.ServeHTTP)

	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "gateway"})
}
