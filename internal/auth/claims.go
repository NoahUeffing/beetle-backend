package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ClaimsData struct {
	UserId uuid.UUID `json:"userId"`

	// TODO: add role/s here so that we don't have to query on each request?
}

type CombinedClaims struct {
	ClaimsData

	jwt.RegisteredClaims
}
