package dto

// RegisterReq represents a user registration request
type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

// RegisterResp represents a user registration response
type RegisterResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// LoginReq represents a user login request
type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResp represents a user login response
type LoginResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshTokenReq represents a refresh token request
type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshTokenResp represents a refresh token response
type RefreshTokenResp struct {
	AccessToken string `json:"accessToken"`
}
