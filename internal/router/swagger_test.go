package router_test

import (
	"beetle/internal/config"
	"beetle/internal/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestAddSwaggerRoutes(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create test config
	cfg := config.Config{
		Logs: config.LogsConfig{
			Level:              log.INFO,
			JSON:               false,
			HideStartupMessage: true,
		},
	}

	// Add Swagger routes
	router.AddSwaggerRoutes(e, cfg)

	// Test cases for different Swagger endpoints
	testCases := []struct {
		name       string
		path       string
		statusCode int
	}{
		{
			name:       "Swagger UI index",
			path:       "/swagger/index.html",
			statusCode: http.StatusOK,
		},
		{
			name:       "Swagger UI assets",
			path:       "/swagger/swagger-ui.css",
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()

			// Serve the request
			e.ServeHTTP(rec, req)

			// Assert the status code
			assert.Equal(t, tc.statusCode, rec.Code)
		})
	}
}
