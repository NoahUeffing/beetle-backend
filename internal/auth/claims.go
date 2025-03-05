package auth

import (
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
