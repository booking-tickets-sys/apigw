package middleware

import (
	"github.com/gin-gonic/gin"
)

// HealthCheck handles health check requests
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "api-gateway",
		})
	}
}
