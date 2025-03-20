package handler

import (
	"beetle/internal/auth"
	"beetle/internal/postgres"
	"beetle/internal/validation"
	"database/sql"
	"errors"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type errorSet struct {
	Code   int
	Errors []error
}

var errorCodes = []errorSet{
	{http.StatusNotFound, []error{
		sql.ErrNoRows,
		postgres.ErrEntityNotFound,
		gorm.ErrRecordNotFound,
	}},
	{http.StatusBadRequest, []error{
		postgres.ErrEntityNonUnique,
		postgres.ErrEntityVersionConflict,
		postgres.ErrInvalidCredentials,
	}},
	{http.StatusUnauthorized, []error{
		auth.ErrUnauthorized,
		echojwt.ErrJWTInvalid,
		echojwt.ErrJWTMissing,
	}},
}

func getCode(err error) (int, bool) {
	for _, c := range errorCodes {
		for _, e := range c.Errors {
			if errors.Is(err, e) {
				return c.Code, true
			}
		}
	}

	return -1, false
}

type Message struct {
	Message string `json:"message"`
}

type FormValidationError struct {
	Message string                             `json:"message"`
	Fields  []validation.TranslatedFieldErrors `json:"fields"`
}

func errorResponse(code int, err error) error {
	return echo.NewHTTPError(code, err.Error())
}

func WrapError(err error) error {
	err = postgres.ConvertErrorIfNeeded(err)
	if code, ok := getCode(err); ok {
		return errorResponse(code, err)
	}
	return err
}
