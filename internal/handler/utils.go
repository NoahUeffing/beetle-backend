package handler

import (
	"beetle/internal/domain"
	"beetle/internal/postgres"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func parseUUID(idStr string) (uuid.UUID, error) {
	if idStr == "" {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, postgres.ErrEntityNotFound
	}
	return id, nil
}

func parseUUIDs(ids []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id) // Trim whitespace!
		if id == "" {
			continue // skip empty strings
		}
		u, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID: %q, error: %w", id, err)
		}
		uuids = append(uuids, u)
	}
	return uuids, nil
}

func buildIDFilter(c Context, param, field string) (*domain.Filter, error) {
	ids := c.QueryParam(param)
	if ids == "" {
		return nil, nil
	}
	idStrings := strings.Split(ids, ",")
	if len(idStrings) > domain.MaxFilterIDs {
		return nil, echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Too many %s IDs", field)})
	}
	parsed, err := parseUUIDs(idStrings)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid %s ID(s)", field)})
	}
	return &domain.Filter{
		Field:    field,
		Operator: "in",
		Value:    parsed,
	}, nil
}

func buildStringFilter(c Context, param, field, operator string) *domain.Filter {
	value := c.QueryParam(param)
	if value == "" {
		return nil
	}
	return &domain.Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}
