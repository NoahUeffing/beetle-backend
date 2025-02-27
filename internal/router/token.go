package router

import (
	"beetle/internal/config"
	"beetle/internal/handler"

	"github.com/labstack/echo/v4"
)

// WithContext wraps a handler function with echo context
func WithContext(h func(handler.Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		if ctx, ok := c.(handler.Context); ok {
			return h(ctx)
		}
		return echo.ErrInternalServerError
	}
}

type tokenRouteProvider struct{}

func (r *tokenRouteProvider) AddPublicRoutes(g *echo.Group, config config.Config) {
	g.POST("/tokens", WithContext(handler.AuthTokenCreate))
}

func (r *tokenRouteProvider) AddPrivateRoutes(g *echo.Group, config config.Config) {
	// No Private Routes
}
