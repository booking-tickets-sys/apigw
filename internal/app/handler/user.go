package handler

import (
	"net/http"

	pb "apigw/client/proto"
	"apigw/internal/app/domains/dto"
	"apigw/internal/app/middleware"
	"apigw/internal/client"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userClient *client.UserServiceClient
	logger     *logrus.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userClient *client.UserServiceClient, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		userClient: userClient,
		logger:     logger,
	}
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorHandler(c, "INVALID_REQUEST", "Invalid request body", h.logger)
		return
	}

	resp, err := h.userClient.Register(c.Request.Context(), &pb.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
	})
	if err != nil {
		middleware.GRPCErrorHandler(c, err, h.logger)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorHandler(c, "INVALID_REQUEST", "Invalid request body", h.logger)
		return
	}

	resp, err := h.userClient.Login(c.Request.Context(), &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		middleware.GRPCErrorHandler(c, err, h.logger)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RefreshToken handles token refresh
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorHandler(c, "INVALID_REQUEST", "Invalid request body", h.logger)
		return
	}

	resp, err := h.userClient.RefreshToken(c.Request.Context(), &pb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		middleware.GRPCErrorHandler(c, err, h.logger)
		return
	}

	c.JSON(http.StatusOK, resp)
}
