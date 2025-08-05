package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// JWTMaker is a JWT token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTTokenMaker creates a new JWT token maker
func NewJWTTokenMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("invalid key size: must be at least 32 characters")
	}
	return &JWTMaker{secretKey: secretKey}, nil
}

// CreateToken creates a new token for a specific user and duration
func (maker *JWTMaker) CreateToken(userID string, duration time.Duration) (string, error) {
	payload, err := NewPayload(userID, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
