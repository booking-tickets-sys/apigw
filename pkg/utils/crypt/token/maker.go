package token

import (
	"github.com/golang-jwt/jwt/v5"
)

// Payload represents the JWT payload
type Payload struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}


