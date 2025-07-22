package handler

import (
	"beetle/internal/postgres"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetCompany godoc
// @Summary Get a company by ID
// @Description Get a companies details from ID
// @Tags company
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} domain.Company
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Company not found"
// @Failure 500 {string} string "Internal server error"
// @Router /company/{id} [get]
func GetCompany(c Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Company ID is required"})
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return postgres.ErrEntityNotFound
	}

	user, err := c.CompanyService.ReadByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, map[string]string{"error": "Company not found"})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve company"})
	}

	return c.JSON(http.StatusOK, user)
}

// GetCompanies godoc
// @Summary Get companies
// @Description Get companies
// @Tags company
// @Accept json
// @Produce json
// @Param name query string false "Name to search"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 12, max: 120)"
// @Success 200 {object} domain.PaginatedResults
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /company [get]
func GetCompanies(c Context) error {
	filter := buildStringFilter(c, "name", "company search", "like")
	companies, err := c.CompanyService.GetCompanies(c.PaginationQuery, *filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve companies"})
	}
	return c.JSON(http.StatusOK, companies)
}
