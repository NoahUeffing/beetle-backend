package router_test

import (
	"beetle/internal/config"
	"beetle/internal/router"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthcheckRouteProvider_AddPublicRoutes(t *testing.T) {
	// Setup
	e := echo.New()
	g := e.Group("/api")
	provider := &router.HealthcheckRouteProvider{}
	cfg := config.Config{}

	// Execute
	provider.AddPublicRoutes(g, cfg)

	// Assert
	routes := e.Routes()
	found := false
	for _, route := range routes {
		if route.Method == "GET" && route.Path == "/api/healthcheck" {
			found = true
			break
		}
	}
	assert.True(t, found, "GET /api/healthcheck route should be registered")
}

func TestHealthcheckRouteProvider_AddPrivateRoutes(t *testing.T) {
	// Setup
	e := echo.New()
	g := e.Group("/api")
	provider := &router.HealthcheckRouteProvider{}
	cfg := config.Config{}

	// Execute
	provider.AddPrivateRoutes(g, cfg)

	// Assert
	routes := e.Routes()
	assert.Empty(t, routes, "No private routes should be registered")
}
