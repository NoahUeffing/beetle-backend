package auth

import (
	"beetle/internal/config"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := config.AuthConfig{Secret: "test-secret"}
	auth := New(cfg)
	assert.NotNil(t, auth)
	assert.Equal(t, cfg, auth.Config)
}

func TestNewToken(t *testing.T) {
	cfg := config.AuthConfig{Secret: "test-secret"}
	auth := New(cfg)

	userId := uuid.New()
	data := &ClaimsData{UserId: userId}

	token, err := auth.NewToken(data)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the token can be parsed and contains correct claims
	parsedToken, err := jwt.ParseWithClaims(token, &CombinedClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(*CombinedClaims)
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
	auth := New(cfg)

	middleware := auth.GetMiddleware()
	assert.NotNil(t, middleware)
}
