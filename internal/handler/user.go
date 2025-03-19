package handler

import (
	"beetle/internal/domain"
	"beetle/internal/validation"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService domain.IUserService
	validator   *validation.Validator
}

func NewUserHandler(userService domain.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validation.New(),
	}
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
func (h *UserHandler) CreateUser(c echo.Context) error {
	var input domain.UserCreateInput

	// Parse request body
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate input using the validator
	if errors := h.validator.Validate(&input); errors != nil {
		translatedErrors := translateErrors(h.validator, &input, errors)

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation failed",
			"errors": translatedErrors,
		})
	}

	// Create user
	user, err := h.userService.CreateUser(&input)
	if err != nil {
		// TODO: Check if the error is a duplicate username or email
		// Log the error for debugging
		log.Printf("Error creating user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
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
func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user"})
	}

	return c.JSON(http.StatusOK, user)
}
