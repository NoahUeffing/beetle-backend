package handler

import (
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
