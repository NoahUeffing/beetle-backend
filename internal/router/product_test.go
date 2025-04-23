package router_test

import (
	"beetle/internal/config"
	"beetle/internal/router"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestProductRouteProvider_AddPublicRoutes(t *testing.T) {
	// Setup
	e := echo.New()
	g := e.Group("/api")
	provider := &router.ProductRouteProvider{}
	cfg := config.Config{}

	// Execute
	provider.AddPublicRoutes(g, cfg)

	// Assert
	routes := e.Routes()
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/product/license/:id"},
		{"GET", "/api/product/licenses"},
		{"GET", "/api/product/dosage-forms"},
		{"GET", "/api/product/dosage-form/:id"},
		{"GET", "/api/product/license/submission-types"},
		{"GET", "/api/product/license/submission-type/:id"},
	}

	foundRoutes := make(map[string]bool)
	for _, route := range routes {
		for _, expected := range expectedRoutes {
			if route.Method == expected.method && route.Path == expected.path {
				foundRoutes[expected.path] = true
			}
		}
	}

	for _, expected := range expectedRoutes {
		assert.True(t, foundRoutes[expected.path], "Route %s %s should be registered", expected.method, expected.path)
	}
}

func TestProductRouteProvider_AddPrivateRoutes(t *testing.T) {
	// Setup
	e := echo.New()
	g := e.Group("/api")
	provider := &router.ProductRouteProvider{}
	cfg := config.Config{}

	// Execute
	provider.AddPrivateRoutes(g, cfg)

	// Assert
	routes := e.Routes()
	assert.Empty(t, routes, "No private routes should be registered")
}
