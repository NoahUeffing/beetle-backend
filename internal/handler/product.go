package handler

import (
	"beetle/internal/postgres"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetProductLicense godoc
// @Summary Get a product license by ID
// @Description Get a product license details from ID
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "Product License ID"
// @Success 200 {object} domain.ProductLicense
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "License not found"
// @Failure 500 {string} string "Internal server error"
// @Router /product/license/{id} [get]
func GetProductLicense(c Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Product License ID is required"})
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return postgres.ErrEntityNotFound
	}

	user, err := c.ProductService.ReadLicenseByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, map[string]string{"error": "License not found"})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve license"})
	}

	return c.JSON(http.StatusOK, user)
}
