package router_test

import (
	"beetle/internal/config"
	"beetle/internal/router"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestTokenRouteProvider_AddPublicRoutes(t *testing.T) {
	// Setup
	e := echo.New()
	g := e.Group("/api")
	provider := &router.TokenRouteProvider{}
	cfg := config.Config{}

	// Execute
	provider.AddPublicRoutes(g, cfg)

	// Assert
	routes := e.Routes()
	found := false
	for _, route := range routes {
		if route.Method == "POST" && route.Path == "/api/tokens" {
			found = true
			break
		}
	}
	assert.True(t, found, "POST /api/tokens route should be registered")
}

func TestTokenRouteProvider_AddPrivateRoutes(t *testing.T) {
	// Setup
	e := echo.New()
	g := e.Group("/api")
	provider := &router.TokenRouteProvider{}
	cfg := config.Config{}

	// Execute
	provider.AddPrivateRoutes(g, cfg)

	// Assert
	routes := e.Routes()
	assert.Empty(t, routes, "No private routes should be registered")
}
