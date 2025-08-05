package handler

import (
	"net/http"

	pb "apigw/client/proto"
	"apigw/internal/app/middleware"
	"apigw/internal/client"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// TicketHandler handles HTTP requests for ticket operations
type OrderHandler struct {
	orderClient *client.OrderServiceClient
	logger      *logrus.Logger
}

// NewTicketHandler creates a new ticket handler
func NewOrderHandler(orderClient *client.OrderServiceClient, logger *logrus.Logger) *OrderHandler {
	return &OrderHandler{
		orderClient: orderClient,
		logger:      logger,
	}
}

// PurchaseTicket handles ticket purchase
func (h *OrderHandler) PurchaseTicket(c *gin.Context) {
	// Get user ID from context (set by JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		middleware.AuthenticationErrorHandler(c, h.logger)
		return
	}

	// Get event ID from URL parameter
	eventID := c.Param("event_id")
	if eventID == "" {
		middleware.ValidationErrorHandler(c, "INVALID_EVENT_ID", "Event ID is required", h.logger)
		return
	}

	resp, err := h.orderClient.PurchaseTicket(c.Request.Context(), &pb.PurchaseRequest{
		EventId: eventID,
		UserId:  userID.(string),
	})
	if err != nil {
		middleware.GRPCErrorHandler(c, err, h.logger)
		return
	}

	c.JSON(http.StatusOK, resp)
}
