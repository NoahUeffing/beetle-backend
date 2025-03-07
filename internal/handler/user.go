package handler

import (
	"beetle/internal/domain"
	"beetle/internal/validation"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
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
		// Get the English translator
		trans, _ := h.validator.Translator.GetTranslator("en")

		// Translate validation errors
		translatedErrors := h.validator.TranslateFormErrors(trans, *errors)

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation failed",
			"errors": translatedErrors,
		})
	}

	// Create user
	user, err := h.userService.CreateUser(&input)
	if err != nil {
		// Log the error for debugging
		log.Printf("Error creating user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	// Return created user
	return c.JSON(http.StatusCreated, user)
}
