package server_test

import (
	"beetle/internal/auth"
	"beetle/internal/config"
	"beetle/internal/handler"
	"beetle/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	// Setup test configuration
	cfg := &config.Config{
		Logs: config.LogsConfig{
			JSON:               false,
			Level:              log.DEBUG,
			HideStartupMessage: false,
		},
	}

	// Setup context configuration
	contextConfig := handler.ContextConfig{
		AuthService: &mockAuthService{},
	}

	// Create new server
	s := server.New(cfg, contextConfig)

	// Assert server was created correctly
	assert.NotNil(t, s)
	assert.NotNil(t, s.Echo)
	assert.Equal(t, cfg, &s.Config)
}

func TestServerErrorHandler(t *testing.T) {
	// Setup test configuration
	cfg := &config.Config{
		Logs: config.LogsConfig{
			JSON:               false,
			Level:              log.DEBUG,
			HideStartupMessage: false,
		},
	}

	// Setup context configuration
	contextConfig := handler.ContextConfig{
		AuthService: &mockAuthService{},
	}

	// Create new server
	s := server.New(cfg, contextConfig)

	// Test error handling
	req := httptest.NewRequest(http.MethodGet, "/not-found", nil)
	rec := httptest.NewRecorder()

	// Make request
	s.Echo.ServeHTTP(rec, req)

	// Assert 404 status code
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// mockAuthService implements the auth.IAuthService interface for testing
type mockAuthService struct{}

func (m *mockAuthService) NewToken(data *auth.ClaimsData) (string, error) {
	return "mock-token", nil
}

func (m *mockAuthService) GetMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}

func (m *mockAuthService) GetClaims(c echo.Context) (*auth.CombinedClaims, error) {
	return &auth.CombinedClaims{}, nil
}
