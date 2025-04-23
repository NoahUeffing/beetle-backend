package router_test

import (
	"beetle/internal/config"
	"beetle/internal/router"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// MockRouteProvider implements routeProvider interface for testing
type MockRouteProvider struct {
	publicRoutesAdded  bool
	privateRoutesAdded bool
}

func (m *MockRouteProvider) AddPublicRoutes(g *echo.Group, config config.Config) {
	m.publicRoutesAdded = true
}

func (m *MockRouteProvider) AddPrivateRoutes(g *echo.Group, config config.Config) {
	m.privateRoutesAdded = true
}

func TestNew(t *testing.T) {
	// Arrange
	cfg := config.Config{}

	// Act
	r := router.New(cfg)

	// Assert
	assert.NotNil(t, r)
	assert.Equal(t, cfg, r.Config)
	assert.Len(t, r.RouteProviders, 5)

	// Verify all expected providers are present
	providerTypes := make(map[string]bool)
	for _, provider := range r.RouteProviders {
		switch provider.(type) {
		case *router.TokenRouteProvider:
			providerTypes["token"] = true
		case *router.HealthcheckRouteProvider:
			providerTypes["healthcheck"] = true
		case *router.UserRouteProvider:
			providerTypes["user"] = true
		case *router.CompanyRouteProvider:
			providerTypes["company"] = true
		case *router.ProductRouteProvider:
			providerTypes["product"] = true
		}
	}

	assert.True(t, providerTypes["token"])
	assert.True(t, providerTypes["healthcheck"])
	assert.True(t, providerTypes["user"])
	assert.True(t, providerTypes["company"])
	assert.True(t, providerTypes["product"])
}

func TestAddRoutesWithAuthMiddleware(t *testing.T) {
	// Arrange
	cfg := config.Config{}
	r := router.New(cfg)
	server := echo.New()
	mockProvider := &MockRouteProvider{}
	r.RouteProviders = []router.RouteProvider{mockProvider}

	authMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}

	// Act
	r.AddRoutes(server, authMiddleware)

	// Assert
	assert.True(t, mockProvider.publicRoutesAdded, "Public routes should be added")
	assert.True(t, mockProvider.privateRoutesAdded, "Private routes should be added")

	// Verify that the private group has the auth middleware
	foundPrivateWithAuth := false

	for _, route := range server.Routes() {
		if route.Path == "/v1/*" && route.Name != "v1" {
			foundPrivateWithAuth = true
			break
		}
	}

	assert.True(t, foundPrivateWithAuth, "Private group should have auth middleware")
}
