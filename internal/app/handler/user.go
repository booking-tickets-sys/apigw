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
	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}).Info("User registration request received")

	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"error":  err.Error(),
		}).Warn("Invalid registration request body")
		middleware.ValidationErrorHandler(c, "INVALID_REQUEST", "Invalid request body", h.logger)
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":   c.Request.Method,
		"path":     c.Request.URL.Path,
		"email":    req.Email,
		"username": req.Username,
	}).Info("Processing user registration")

	resp, err := h.userClient.Register(c.Request.Context(), &pb.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
	})
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"email":  req.Email,
			"error":  err.Error(),
		}).Error("User registration failed")
		middleware.GRPCErrorHandler(c, err, h.logger)
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"email":  req.Email,
	}).Info("User registration successful")

	c.JSON(http.StatusCreated, dto.RegisterResp{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	})
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}).Info("User login request received")

	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"error":  err.Error(),
		}).Warn("Invalid login request body")
		middleware.ValidationErrorHandler(c, "INVALID_REQUEST", "Invalid request body", h.logger)
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"email":  req.Email,
	}).Info("Processing user login")

	resp, err := h.userClient.Login(c.Request.Context(), &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"email":  req.Email,
			"error":  err.Error(),
		}).Error("User login failed")
		middleware.GRPCErrorHandler(c, err, h.logger)
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"email":  req.Email,
	}).Info("User login successful")

	c.JSON(http.StatusOK, dto.LoginResp{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	})
}

// RefreshToken handles token refresh
func (h *UserHandler) RefreshToken(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}).Info("Token refresh request received")

	var req dto.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"error":  err.Error(),
		}).Warn("Invalid refresh token request body")
		middleware.ValidationErrorHandler(c, "INVALID_REQUEST", "Invalid request body", h.logger)
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	}).Info("Processing token refresh")

	resp, err := h.userClient.RefreshToken(c.Request.Context(), &pb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"error":  err.Error(),
		}).Error("Token refresh failed")
		middleware.GRPCErrorHandler(c, err, h.logger)
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	}).Info("Token refresh successful")

	c.JSON(http.StatusOK, dto.RefreshTokenResp{
		AccessToken: resp.AccessToken,
	})
}
