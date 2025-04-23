package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ClaimsData struct {
	UserId uuid.UUID `json:"userId"`
}

type CombinedClaims struct {
	ClaimsData

	jwt.RegisteredClaims
}

// Valid implements the jwt.Claims interface
func (c *CombinedClaims) Valid() error {
	now := time.Now().UTC()

	if c.ExpiresAt != nil && now.After(c.ExpiresAt.Time) {
		return jwt.ErrTokenExpired
	}

	if c.NotBefore != nil && now.Before(c.NotBefore.Time) {
		return jwt.ErrTokenNotValidYet
	}

	return nil
}
