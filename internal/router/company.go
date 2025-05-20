package router

import (
	"beetle/internal/config"
	"beetle/internal/handler"
	"beetle/internal/router/middleware"

	"github.com/labstack/echo/v4"
)

type CompanyRouteProvider struct{}

func (r *CompanyRouteProvider) AddPublicRoutes(g *echo.Group, config config.Config) {
	// Public routes
	g.GET("/company/:id", WithContext(handler.GetCompany))
	g.GET("/company", WithContext(handler.GetCompanies, middleware.Paginate))
}

func (r *CompanyRouteProvider) AddPrivateRoutes(g *echo.Group, config config.Config) {
	// Private routes (requiring authentication)
}
