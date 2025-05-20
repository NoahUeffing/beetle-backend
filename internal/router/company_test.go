package router_test

import (
	"beetle/internal/config"
	"beetle/internal/router"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCompanyRouteProvider_AddPublicRoutes(t *testing.T) {
	// Setup
	e := echo.New()
	g := e.Group("/api")
	provider := &router.CompanyRouteProvider{}
	cfg := config.Config{}

	// Execute
	provider.AddPublicRoutes(g, cfg)

	// Assert
	routes := e.Routes()
	found := false
	for _, route := range routes {
		if route.Method == "GET" && route.Path == "/api/company/:id" {
			found = true
			break
		}
	}
	assert.True(t, found, "GET /api/company/:id route should be registered")
	assert.True(t, found, "GET /api/company route should be registered")
}

func TestCompanyRouteProvider_AddPrivateRoutes(t *testing.T) {
	// Setup
	e := echo.New()
	g := e.Group("/api")
	provider := &router.CompanyRouteProvider{}
	cfg := config.Config{}

	// Execute
	provider.AddPrivateRoutes(g, cfg)

	// Assert
	routes := e.Routes()
	assert.Empty(t, routes, "No private routes should be registered")
}
