package dto

// RegisterReq represents a user registration request
type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginReq represents a user login request
type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenReq represents a refresh token request
type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
