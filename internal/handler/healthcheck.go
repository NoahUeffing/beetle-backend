package handler

import (
	"beetle/internal/healthcheck"
	"net/http"
)

// HealthCheck godoc
// @Summary Allows healthcheck of server in general and substatuses of various internal services
// @ID v1-healthcheck
// @Tags test
// @Produce json
// @Success 200 {object} healthcheck.Status
// @Failure 503 {object} healthcheck.Status
// @Router /v1/healthcheck [get]
func HealthCheck(c Context) error {
	result := []healthcheck.Status{}
	allUp := true
	var status healthcheck.Status

	for _, hcs := range c.HealthCheckServices {
		status = hcs.Check()
		result = append(result, status)
		if !status.Up {
			allUp = false
		}
	}

	if !allUp {
		return c.JSON(http.StatusServiceUnavailable, result)
	}

	return c.JSON(http.StatusOK, result)
}
