package handler

import (
	"net/http"
	"time"

	pb "apigw/api/proto"

	"github.com/gin-gonic/gin"

	"apigw/internal/client"
)

type UserHandler struct {
	userClient *client.UserServiceClient
}

func NewUserHandler(userClient *client.UserServiceClient) *UserHandler {
	return &UserHandler{
		userClient: userClient,
	}
}

// RegisterRequest represents the HTTP request for user registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the HTTP request for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents the HTTP request for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UserResponse represents the HTTP response for user data
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// AuthResponse represents the HTTP response for authentication operations
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
}

// TokenResponse represents the HTTP response for token operations
type TokenResponse struct {
	AccessToken string `json:"accessToken"`
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &pb.RegisterRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	resp, err := h.userClient.Register(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResp := UserResponse{
		ID:        resp.User.Id,
		Email:     resp.User.Email,
		Username:  resp.User.Username,
		CreatedAt: resp.User.CreatedAt.AsTime(),
		UpdatedAt: resp.User.UpdatedAt.AsTime(),
	}

	authResp := AuthResponse{
		User:         userResp,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}

	c.JSON(http.StatusCreated, authResp)
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.userClient.Login(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userResp := UserResponse{
		ID:        resp.User.Id,
		Email:     resp.User.Email,
		Username:  resp.User.Username,
		CreatedAt: resp.User.CreatedAt.AsTime(),
		UpdatedAt: resp.User.UpdatedAt.AsTime(),
	}

	authResp := AuthResponse{
		User:         userResp,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}

	c.JSON(http.StatusOK, authResp)
}

// RefreshToken handles token refresh
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &pb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	resp, err := h.userClient.RefreshToken(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokenResp := TokenResponse{
		AccessToken: resp.AccessToken,
	}

	c.JSON(http.StatusOK, tokenResp)
}
