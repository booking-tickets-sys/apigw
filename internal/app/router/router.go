package router

import (
	"apigw/internal/app/config"
	"apigw/internal/app/handler"
	"apigw/internal/app/middleware"
	"apigw/internal/client"
	"apigw/pkg/utils/crypt/token"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRouter configures and returns the HTTP router
func SetupRouter(
	cfg *config.Config,
	userClient *client.UserServiceClient,
	orderClient *client.OrderServiceClient,
	jwtMaker *token.JWTMaker,
	logger *logrus.Logger,
) *gin.Engine {
	// Set Gin mode
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware(logger))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"service":   cfg.App.Name,
			"version":   cfg.App.Version,
			"timestamp": "2024-01-01T00:00:00Z",
		})
	})

	// Create handlers
	userHandler := handler.NewUserHandler(userClient, logger)
	orderHandler := handler.NewOrderHandler(orderClient, logger)

	// Create JWT middleware
	jwtMiddleware := middleware.JWTMiddleware(jwtMaker, logger)

	// API routes
	api := router.Group("/api/v1")
	{
		// User routes (no authentication required)
		users := api.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
			users.POST("/refresh", userHandler.RefreshToken)
		}

		// Ticket routes (authentication required)
		tickets := api.Group("/tickets")
		tickets.Use(jwtMiddleware)
		{
			tickets.POST("/:event_id/purchase", orderHandler.PurchaseTicket)
		}
	}

	return router
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
