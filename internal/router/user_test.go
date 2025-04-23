package router_test

import (
	"beetle/internal/config"
	"beetle/internal/router"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserRouteProvider_AddPublicRoutes(t *testing.T) {
	// Setup
	e := echo.New()
	g := e.Group("/api")
	provider := &router.UserRouteProvider{}
	cfg := config.Config{}

	// Execute
	provider.AddPublicRoutes(g, cfg)

	// Assert
	routes := e.Routes()
	found := false
	for _, route := range routes {
		if route.Method == "POST" && route.Path == "/api/user" {
			found = true
			break
		}
	}
	assert.True(t, found, "POST /api/user route should be registered")
}

func TestUserRouteProvider_AddPrivateRoutes(t *testing.T) {
	// Setup
	e := echo.New()
	g := e.Group("/api")
	provider := &router.UserRouteProvider{}
	cfg := config.Config{}

	// Execute
	provider.AddPrivateRoutes(g, cfg)

	// Assert
	routes := e.Routes()
	found := false
	for _, route := range routes {
		if route.Method == "GET" && route.Path == "/api/user/:id" {
			found = true
			break
		}
	}
	assert.True(t, found, "GET /api/user/:id route should be registered")
}
