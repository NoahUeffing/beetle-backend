package handler

import (
	"beetle/internal/postgres"
	"net/http"

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
