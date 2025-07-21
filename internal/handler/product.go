package handler

import (
	"beetle/internal/domain"
	"net/http"
	"strings"

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
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		return err
	}

	license, err := c.ProductService.ReadLicenseByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, map[string]string{"error": "License not found"})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve license"})
	}

	return c.JSON(http.StatusOK, license)
}

// GetLicenses godoc
// @Summary Get product licenses
// @Description Get product licenses
// @Tags product
// @Accept json
// @Produce json
// @Param companies query string false "Comma-separated list of IDs, max of 5"
// @Param forms query string false "Comma-separated list of IDs, max of 5"
// @Param name query string false "Name to search"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 12, max: 120)"
// @Success 200 {object} domain.PaginatedResults
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /product/licenses [get]
func GetLicenses(c Context) error {
	// TODO: Abstract some logic
	var filters []domain.Filter

	// Company IDs
	if ids := c.QueryParam("companies"); ids != "" {
		idStrings := strings.Split(ids, ",")
		if len(idStrings) > 5 { // limit
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "Too many company IDs"})
		}
		parsed, err := parseUUIDs(idStrings)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "Invalid company ID(s)"})
		}
		filters = append(filters, domain.Filter{
			Field:    "company",
			Operator: "in",
			Value:    parsed,
		})
	}
	// Dosage Form IDs
	if ids := c.QueryParam("forms"); ids != "" {
		idStrings := strings.Split(ids, ",")
		if len(idStrings) > 20 { // limit
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "Too many form IDs"})
		}
		parsed, err := parseUUIDs(idStrings)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "Invalid dosage form ID(s)"})
		}
		filters = append(filters, domain.Filter{
			Field:    "form",
			Operator: "in",
			Value:    parsed,
		})
	}

	// Name search
	if name := c.QueryParam("name"); name != "" {
		filters = append(filters, domain.Filter{
			Field:    "name",
			Operator: "like",
			Value:    name,
		})
	}
	// Pass filters to service
	licenses, err := c.ProductService.GetLicenses(c.PaginationQuery, filters...)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, licenses)
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

// GetDosageFormByID godoc
// @Summary Get a dosage form by ID
// @Description Get a dosage form details from ID
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "Dosage Form ID"
// @Success 200 {object} domain.DosageForm
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /product/dosage-form/{id} [get]
func GetDosageFormByID(c Context) error {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		return err
	}
	dosageForm, err := c.ProductService.ReadDosageFormByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve dosage form"})
	}
	return c.JSON(http.StatusOK, dosageForm)
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

// GetSubmissionTypeByID godoc
// @Summary Get a submission type by ID
// @Description Get a submission type details from ID
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "Submission Type ID"
// @Success 200 {object} domain.SubmissionType
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /product/license/submission-type/{id} [get]
func GetSubmissionTypeByID(c Context) error {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		return err
	}
	submissionType, err := c.ProductService.ReadSubmissionTypeByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve submission type"})
	}
	return c.JSON(http.StatusOK, submissionType)
}
