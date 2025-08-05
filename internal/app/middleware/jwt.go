package middleware

import (
	"apigw/pkg/utils/crypt/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// JWTMiddleware creates JWT authentication middleware
func JWTMiddleware(
	jwtMaker *token.JWTMaker,
	logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for certain paths
		if shouldSkipAuth(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Error("Authorization header missing")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "AUTHENTICATION_ERROR",
				"code":    "MISSING_TOKEN",
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if token starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Error("Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "AUTHENTICATION_ERROR",
				"code":    "INVALID_TOKEN_FORMAT",
				"message": "Token must be in format: Bearer <token>",
			})
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		user, err := jwtMaker.VerifyToken(token)
		if err != nil {
			logger.WithError(err).Error("Token validation failed")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "AUTHENTICATION_ERROR",
				"code":    "INVALID_TOKEN",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", user.UserID)

		c.Next()
	}
}

// shouldSkipAuth checks if authentication should be skipped for the given path
func shouldSkipAuth(path string) bool {
	skipPaths := []string{
		"/health",
		"/api/v1/users/register",
		"/api/v1/users/login",
		"/api/v1/users/refresh",
	}

	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}

	return false
}
