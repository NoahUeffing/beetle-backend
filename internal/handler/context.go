package handler

import (
	"beetle/internal/auth"
	"beetle/internal/domain"
	"beetle/internal/healthcheck"
	"beetle/internal/validation"

	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	ContextConfig
	User            *domain.User
	PaginationQuery *domain.PaginationQuery
}

type ContextConfig struct {
	AuthService         auth.IAuthService
	UserService         domain.IUserService
	CompanyService      domain.ICompanyService
	ValidationService   validation.Validator
	HealthCheckServices []healthcheck.IHealthCheckService
	ProductService      domain.IProductService
}
