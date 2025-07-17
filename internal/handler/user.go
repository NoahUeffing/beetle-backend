package handler

import (
	"beetle/internal/domain"
	"beetle/internal/postgres"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

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
// @Router /user [post]
func CreateUser(c Context) error {
	var input domain.UserCreateInput

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if errors := c.validate(&input); errors != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors)
	}

	user, err := c.UserService.CreateUser(&input)
	if err == postgres.ErrUsernameTaken || err == postgres.ErrEmailAlreadyAssociated {
		log.Printf("Error creating user: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	token, err := c.UserService.CreateAuthToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &UserCreateResponse{
		User:  *user,
		Token: token.Token,
	})
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a user's details by their ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} domain.User
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /user/{id} [get]
// @Security JWTToken
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

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user's details
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body domain.User true "User object"
// @Success 200 {object} domain.User
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /user/{id} [put]
// @Security JWTToken
func UpdateUser(c Context) error {
	u := &domain.User{}
	if err := c.Bind(u); err != nil {
		return err
	}

	if formErrs := c.validate(u); formErrs != nil {
		return echo.NewHTTPError(http.StatusBadRequest, formErrs)
	}

	u, err := c.UserService.Update(u)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}

// UserDelete godoc
// @Summary Delete a User account
// @Description Delete a user and remove personal info from db.
// @ID v1-user-delete
// @Tags user
// @Produce json
// @Param PasswordInput body domain.PasswordInput true "Authenticated users password"
// @Success 200 {object} handler.Message
// @Failure 400 {object} handler.FormValidationError
// @Failure 403 {object} handler.Message
// @Failure 500 {object} handler.Message
// @Router /user/delete [put]
// @Security JWTToken
func UserDelete(c Context) error {
	password := &domain.PasswordInput{}
	if err := c.Bind(password); err != nil {
		return err
	}

	if formErrs := c.validate(password); formErrs != nil {
		return echo.NewHTTPError(http.StatusBadRequest, formErrs)
	}

	if err := c.UserService.CheckPassword(c.User, password.Password); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err)
	}

	if err := c.UserService.Delete(c.User); err != nil {
		return err
	}
	c.Response().Header().Set(domain.UserHeaderAuthentication, domain.UserHeaderAuthenticatedFalse)
	return c.JSON(http.StatusOK, Message{Message: "User deleted"})
}

// UserResetPasswordCreate godoc
// @Summary Create a new email confirmation for forgotten password
// @Description Creates a new email confirmation to trigger reset password
// @ID v1-user-reset-password-code
// @Tags user
// @Produce json
// @Param emailInput body domain.EmailInput true "Password reset form input"
// @Success 200 {object} handler.Message
// @Failure 400 {object} handler.Message
// @Failure 500 {object} handler.Message
// @Router /user/password-reset-codes [post]
func UserResetPasswordCreate(c Context) error {
	email := &domain.EmailInput{}
	if err := c.Bind(email); err != nil {
		return err
	}

	if formErrs := c.validate(email); formErrs != nil {
		return echo.NewHTTPError(http.StatusBadRequest, formErrs)
	}

	user, err := c.UserService.ReadByEmail(email.Email)
	fmt.Println("Point 1")

	switch err {
	case nil:
		fmt.Println("Point 2")
		err = c.UserService.ResetPasswordCreate(user)
		if err != nil {
			c.Echo().Logger.Errorf("Error sending confirmation email: %e\n", err)
			return err
		}
	case sql.ErrNoRows, postgres.ErrEntityNotFound:
		fmt.Println("Point 3")
		return c.JSON(http.StatusOK, Message{Message: "Confirmation sent to " + email.Email})
	default:
		fmt.Println("Point 4")
		c.Echo().Logger.Errorf("Error sending confirmation email: %e\n", err)
		return err
	}

	return c.JSON(http.StatusOK, Message{Message: "Email confirmation sent to " + email.Email})
}

// UserResetPasswordConfirm godoc
// @Summary Confirm a password reset code
// @Description Confirms a password reset
// @ID v1-user-reset-password-confirm
// @Tags user
// @Produce json
// @Param member body domain.PasswordResetInput true "Password reset confirm input"
// @Success 200 {object} handler.Message
// @Failure 400 {object} handler.Message
// @Failure 500 {object} handler.Message
// @Router /user/password-reset-confirm [post]
func UserResetPasswordConfirm(c Context) error {
	pass := &domain.PasswordResetInput{}
	if err := c.Bind(pass); err != nil {
		return err
	}

	if formErrs := c.validate(pass); formErrs != nil {
		return echo.NewHTTPError(http.StatusBadRequest, formErrs)
	}

	if err := c.UserService.ResetPasswordConfirm(pass); err != nil {
		if err == postgres.ErrPasswordResetCodeInvalid {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, Message{Message: "Password reset confirmed"})
}
