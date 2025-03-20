package auth

import (
	"errors"
	"time"

	"beetle/internal/config"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const tokenDurationDays = 14
const dayLengthHours = 24
const tokenDuration = tokenDurationDays * dayLengthHours * time.Hour

var ErrUnauthorized = errors.New("invalid or missing authorization token")

type Auth struct {
	Config config.AuthConfig
}

type IAuthService interface {
	NewToken(data *ClaimsData) (string, error)
	GetMiddleware() echo.MiddlewareFunc
	GetClaims(c echo.Context) (*CombinedClaims, error)
}

func New(c config.AuthConfig) *Auth {
	return &Auth{Config: c}
}

func (a *Auth) NewToken(data *ClaimsData) (string, error) {
	return a.buildToken(&CombinedClaims{
		ClaimsData: *data,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(tokenDuration)),
		},
	})
}

func (a *Auth) buildToken(claims *CombinedClaims) (string, error) {
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	return token.SignedString([]byte(a.Config.Secret))
}

func (a *Auth) GetMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(a.Config.Secret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims { return &CombinedClaims{} },
	})
}

func (a *Auth) GetClaims(c echo.Context) (*CombinedClaims, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok || token == nil {
		return nil, ErrUnauthorized
	}
	claims := token.Claims.(*CombinedClaims)
	if claims == nil {
		return nil, ErrUnauthorized
	}
	return claims, nil
}
