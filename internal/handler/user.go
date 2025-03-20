package handler

import (
	"beetle/internal/domain"
	"beetle/internal/postgres"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// TODO: Implement
type UserCreateResponse struct {
	domain.User
	Token string `json:"token"`
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags user
// @Accept json
// @Produce json
// @Param user body domain.UserCreateInput true "User object"
// @Success 201 {object} domain.User
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/user [post]
func CreateUser(c Context) error {
	var input domain.UserCreateInput

	// Parse request body
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate input using the validator
	if errors := c.validate(&input); errors != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors)
	}

	// Create user
	user, err := c.UserService.CreateUser(&input)
	if err != nil {
		// TODO: Check if the error is a duplicate username or email
		// Log the error for debugging
		log.Printf("Error creating user: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Return created user
	return c.JSON(http.StatusCreated, user)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a user's details by their ID
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} domain.User
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/user/{id} [get]
func GetUser(c Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return postgres.ErrEntityNotFound
	}

	user, err := c.UserService.ReadByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user"})
	}

	return c.JSON(http.StatusOK, user)
}
