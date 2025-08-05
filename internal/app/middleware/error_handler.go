package middleware

import (
	"net/http"

	"apigw/internal/app/domains/errs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorHandlerMiddleware provides centralized error handling for gRPC errors
func ErrorHandlerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are any errors in the context
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Convert gRPC error to HTTP error
			httpErr := errs.GRPCToHTTPError(err)

			logger.WithError(err).WithFields(logrus.Fields{
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"status":     httpErr.Status,
				"error_code": httpErr.Code,
			}).Error("Request failed")

			c.JSON(httpErr.Status, httpErr)
			return
		}
	}
}

// GRPCErrorHandler is a helper function that can be used in handlers to handle gRPC errors
func GRPCErrorHandler(c *gin.Context, err error, logger *logrus.Logger) {
	if err == nil {
		return
	}

	// Convert gRPC error to HTTP error
	httpErr := errs.GRPCToHTTPError(err)

	logger.WithError(err).WithFields(logrus.Fields{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"status":     httpErr.Status,
		"error_code": httpErr.Code,
		"grpc_code":  errs.GetGRPCCode(err).String(),
	}).Error("gRPC call failed")

	c.JSON(httpErr.Status, httpErr)
}

// ValidationErrorHandler handles validation errors
func ValidationErrorHandler(c *gin.Context, code, message string, logger *logrus.Logger) {
	httpErr := errs.NewHTTPError("VALIDATION_ERROR", code, message, http.StatusBadRequest)

	logger.WithFields(logrus.Fields{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"error_code": code,
	}).Warn("Validation error")

	c.JSON(httpErr.Status, httpErr)
}

// AuthenticationErrorHandler handles authentication errors
func AuthenticationErrorHandler(c *gin.Context, logger *logrus.Logger) {
	httpErr := errs.ErrUnauthorized

	logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	}).Warn("Authentication failed")

	c.JSON(httpErr.Status, httpErr)
}
