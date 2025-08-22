package auth

import (
	"errors"
	"time"

	"beetle/internal/config"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const dayLengthHours = 24
const authTokenDuration = 14 * dayLengthHours * time.Hour
const refreshTokenDuration = 30 * dayLengthHours * time.Hour

var ErrUnauthorized = errors.New("invalid or missing authorization token")

type Tokens struct {
	AccessToken  string
	RefreshToken string
	// TODO: Should this be stored in database, and only sent back in http instead of json?
}

type Auth struct {
	Config config.AuthConfig
}

type IAuthService interface {
	NewTokens(data *ClaimsData) (*Tokens, error)
	GetMiddleware() echo.MiddlewareFunc
	GetClaims(c echo.Context) (*CombinedClaims, error)
}

func New(c config.AuthConfig) *Auth {
	return &Auth{Config: c}
}

func (a *Auth) NewTokens(data *ClaimsData) (*Tokens, error) {
	accessToken, err := a.buildToken(&CombinedClaims{
		ClaimsData: *data,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(authTokenDuration)),
		},
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.buildToken(&CombinedClaims{
		ClaimsData: *data,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(refreshTokenDuration)),
		},
	})
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *Auth) buildToken(claims *CombinedClaims) (string, error) {
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	return token.SignedString([]byte(a.Config.Secret))
}

func authSkipper(c echo.Context) bool {
	return middleware.DefaultSkipper(c)
}

func (a *Auth) GetMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		Skipper:       authSkipper,
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
