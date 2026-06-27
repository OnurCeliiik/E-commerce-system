package routes

import (
	"github.com/OnurCeliiik/ecommerce/services/user/handlers"
	"github.com/OnurCeliiik/ecommerce/services/user/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB                    *gorm.DB
	UserHandler           *handlers.UserHandler
	TokenValidator        middleware.TokenValidator
	InternalAuthMiddleware gin.HandlerFunc
}

func RegisterRoutes(router *gin.Engine, deps Dependencies) {
	router.GET("/health", handlers.HealthCheck(deps.DB))

	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", deps.UserHandler.Register)
		v1.POST("/login", deps.UserHandler.Login)
		v1.GET("/me", middleware.Auth(deps.TokenValidator), deps.UserHandler.Me)

		internal := v1.Group("/internal")
		internal.Use(deps.InternalAuthMiddleware)
		internal.GET("/users/:id", deps.UserHandler.GetUserEmailInternal)
	}
}
