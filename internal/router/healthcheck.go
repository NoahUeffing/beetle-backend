package router

import (
	"beetle/internal/config"
	"beetle/internal/handler"

	"github.com/labstack/echo/v4"
)

type HealthcheckRouteProvider struct{}

func (r *HealthcheckRouteProvider) AddPublicRoutes(g *echo.Group, config config.Config) {
	g.GET("/healthcheck", WithContext(handler.HealthCheck))
}

func (r *HealthcheckRouteProvider) AddPrivateRoutes(g *echo.Group, config config.Config) {
	// No private routes
}
