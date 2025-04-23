package auth_test

import (
	"beetle/internal/auth"
	"beetle/internal/config"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const tokenDuration = 14 * 24 * time.Hour

func TestNew(t *testing.T) {
	cfg := config.AuthConfig{Secret: "test-secret"}
	a := auth.New(cfg)
	assert.NotNil(t, a)
	assert.Equal(t, cfg, a.Config)
}

func TestNewToken(t *testing.T) {
	cfg := config.AuthConfig{Secret: "test-secret"}
	a := auth.New(cfg)

	userId := uuid.New()
	data := &auth.ClaimsData{UserId: userId}

	token, err := a.NewToken(data)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the token can be parsed and contains correct claims
	parsedToken, err := jwt.ParseWithClaims(token, &auth.CombinedClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(*auth.CombinedClaims)
	assert.True(t, ok)
	assert.Equal(t, userId, claims.UserId)

	// Verify expiration time
	exp := claims.ExpiresAt.Time
	expectedExp := time.Now().UTC().Add(tokenDuration)
	assert.True(t, exp.After(time.Now().UTC()))
	assert.True(t, exp.Before(expectedExp.Add(time.Hour))) // Allow 1 hour margin
}

func TestGetMiddleware(t *testing.T) {
	cfg := config.AuthConfig{Secret: "test-secret"}
	a := auth.New(cfg)

	middleware := a.GetMiddleware()
	assert.NotNil(t, middleware)
}
