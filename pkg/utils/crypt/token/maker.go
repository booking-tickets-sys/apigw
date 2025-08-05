package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Payload represents the JWT payload
type Payload struct {
	UserID    string    `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	jwt.RegisteredClaims
}

// NewPayload creates a new token payload
func NewPayload(userID string, duration time.Duration) (*Payload, error) {
	now := time.Now()
	payload := &Payload{
		UserID:    userID,
		IssuedAt:  now,
		ExpiredAt: now.Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is valid
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
