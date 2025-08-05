package handler

import (
	"net/http"

	pb "apigw/client/proto"
	"apigw/internal/app/middleware"
	"apigw/internal/client"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// OrderHandler handles HTTP requests for order operations
type OrderHandler struct {
	orderClient *client.OrderServiceClient
	logger      *logrus.Logger
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderClient *client.OrderServiceClient, logger *logrus.Logger) *OrderHandler {
	return &OrderHandler{
		orderClient: orderClient,
		logger:      logger,
	}
}

// PurchaseTicket handles ticket purchase
func (h *OrderHandler) PurchaseTicket(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}).Info("Ticket purchase request received")

	// Get user ID from context (set by JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Warn("Authentication failed - user_id not found in context")
		middleware.AuthenticationErrorHandler(c, h.logger)
		return
	}

	// Get event ID from URL parameter
	eventID := c.Param("event_id")
	if eventID == "" {
		h.logger.WithFields(logrus.Fields{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"user_id": userID,
		}).Warn("Invalid event ID - event_id parameter is empty")
		middleware.ValidationErrorHandler(c, "INVALID_EVENT_ID", "Event ID is required", h.logger)
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":   c.Request.Method,
		"path":     c.Request.URL.Path,
		"user_id":  userID,
		"event_id": eventID,
	}).Info("Processing ticket purchase")

	resp, err := h.orderClient.PurchaseTicket(c.Request.Context(), &pb.PurchaseRequest{
		EventId: eventID,
		UserId:  userID.(string),
	})
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"user_id":  userID,
			"event_id": eventID,
			"error":    err.Error(),
		}).Error("Ticket purchase failed")
		middleware.GRPCErrorHandler(c, err, h.logger)
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":   c.Request.Method,
		"path":     c.Request.URL.Path,
		"user_id":  userID,
		"event_id": eventID,
		"status":   resp.Status,
	}).Info("Ticket purchase successful")

	c.JSON(http.StatusOK, resp)
}
