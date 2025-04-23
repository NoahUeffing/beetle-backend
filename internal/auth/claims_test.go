package auth_test

import (
	"testing"
	"time"

	"beetle/internal/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCombinedClaims_Valid(t *testing.T) {
	now := time.Now().UTC()
	userId := uuid.New()

	tests := []struct {
		name    string
		claims  *auth.CombinedClaims
		wantErr error
	}{
		{
			name: "valid token",
			claims: &auth.CombinedClaims{
				ClaimsData: auth.ClaimsData{
					UserId: userId,
				},
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
					NotBefore: jwt.NewNumericDate(now.Add(-time.Hour)),
				},
			},
			wantErr: nil,
		},
		{
			name: "expired token",
			claims: &auth.CombinedClaims{
				ClaimsData: auth.ClaimsData{
					UserId: userId,
				},
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour)),
				},
			},
			wantErr: jwt.ErrTokenExpired,
		},
		{
			name: "token not yet valid",
			claims: &auth.CombinedClaims{
				ClaimsData: auth.ClaimsData{
					UserId: userId,
				},
				RegisteredClaims: jwt.RegisteredClaims{
					NotBefore: jwt.NewNumericDate(now.Add(time.Hour)),
				},
			},
			wantErr: jwt.ErrTokenNotValidYet,
		},
		{
			name: "token with no time constraints",
			claims: &auth.CombinedClaims{
				ClaimsData: auth.ClaimsData{
					UserId: userId,
				},
				RegisteredClaims: jwt.RegisteredClaims{},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.claims.Valid()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
