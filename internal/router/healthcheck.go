package router

import (
	"beetle/internal/config"
	"beetle/internal/handler"

	"github.com/labstack/echo/v4"
)

type healthcheckRouteProvider struct{}

func (r *healthcheckRouteProvider) AddPublicRoutes(g *echo.Group, config config.Config) {
	g.GET("/healthcheck", handler.HealthCheck)
}

func (r *healthcheckRouteProvider) AddPrivateRoutes(g *echo.Group, config config.Config) {
	// No private routes
}
