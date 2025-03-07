package handler

import (
	"beetle/internal/healthcheck"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck godoc
// @Summary Allows healthcheck of server in general and substatuses of various internal services
// @ID v1-healthcheck
// @Tags test
// @Produce json
// @Success 200 {object} healthcheck.Status
// @Failure 503 {object} healthcheck.Status
// @Router /v1/healthcheck [get]
func HealthCheck(c echo.Context) error {
	// For now, just return a simple health check response
	status := healthcheck.Status{
		Name: "API",
		Up:   true,
		Messages: []string{
			"API is healthy",
		},
	}

	return c.JSON(http.StatusOK, []healthcheck.Status{status})
}
