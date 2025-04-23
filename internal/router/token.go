package router

import (
	"beetle/internal/config"
	"beetle/internal/handler"

	"github.com/labstack/echo/v4"
)

type TokenRouteProvider struct{}

func (r *TokenRouteProvider) AddPublicRoutes(g *echo.Group, config config.Config) {
	g.POST("/tokens", WithContext(handler.AuthTokenCreate))
}

func (r *TokenRouteProvider) AddPrivateRoutes(g *echo.Group, config config.Config) {
	// No Private Routes
}
