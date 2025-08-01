package router

import (
	"apigw/internal/handler"
	"apigw/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", middleware.HealthCheck())

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// User routes
		user := v1.Group("/users")
		{
			user.POST("/register", userHandler.Register)
			user.POST("/login", userHandler.Login)
			user.POST("/refresh", userHandler.RefreshToken)
		}

	}

	return router
}
