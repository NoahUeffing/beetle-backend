package handler

import (
	"beetle/internal/domain"
	"beetle/internal/postgres"
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
)

func GetHeader(c Context, key string) string {
	return c.Request().Header.Get(key)
}

// AuthTokenCreate godoc
// @Summary Get a new authorization token via email and password
// @Description Creates a new JWT auth token bearing the user's identity, which should be used to authorize further requests.
// @ID v1-authtoken-create
// @Tags auth
// @Produce json
// @Param credentials body domain.UserAuthInput true "Login form input"
// @Success 200 {object} domain.User
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/tokens [post]
func AuthTokenCreate(c Context) error {
	i := &domain.UserAuthInput{}
	if err := c.Bind(i); err != nil {
		return err
	}

	if formErrs := c.Validate(i); formErrs != nil {
		return echo.NewHTTPError(http.StatusBadRequest, formErrs)
	}

	m, err := c.UserService.ReadByEmail(i.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Logger().Debug("Login error: no user with email ", i.Email)
			err = postgres.ErrInvalidCredentials
		}
		return err
	}

	if err = c.UserService.CheckPassword(m, i.Password); err != nil {
		c.Logger().Debug("Login error: password incorrect for user ", i.Email)
		return err
	}

	t, err := c.UserService.CreateAuthToken(m)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}
