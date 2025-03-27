package router

import (
	"beetle/internal/config"
	"beetle/internal/handler"

	"github.com/labstack/echo/v4"
)

type tokenRouteProvider struct{}

func (r *tokenRouteProvider) AddPublicRoutes(g *echo.Group, config config.Config) {
	g.POST("/tokens", WithContext(handler.AuthTokenCreate))
}

func (r *tokenRouteProvider) AddPrivateRoutes(g *echo.Group, config config.Config) {
	// No Private Routes
}
