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

// GetDosageForms godoc
// @Summary Get all possible dosage forms for a product
// @Description Get all possible dosage forms for a product license
// @Tags product
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 12, max: 120)"
// @Success 200 {object} domain.PaginatedResults
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /product/dosage-forms [get]
func GetDosageForms(c Context) error {
	dosageForms, err := c.ProductService.GetDosageForms(c.PaginationQuery)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve dosage forms"})
	}
	return c.JSON(http.StatusOK, dosageForms)
}

// GetSubmissionTypes godoc
// @Summary Get all possible submission types for a product license
// @Description Get all possible submission types for a product license
// @Tags product
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 12, max: 120)"
// @Success 200 {object} domain.PaginatedResults
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /product/license/submission-types [get]
func GetSubmissionTypes(c Context) error {
	submissionTypes, err := c.ProductService.GetSubmissionTypes(c.PaginationQuery)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve submission types"})
	}
	return c.JSON(http.StatusOK, submissionTypes)
}
