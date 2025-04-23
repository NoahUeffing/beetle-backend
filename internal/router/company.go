package router

import (
	"beetle/internal/config"
	"beetle/internal/handler"

	"github.com/labstack/echo/v4"
)

type CompanyRouteProvider struct{}

func (r *CompanyRouteProvider) AddPublicRoutes(g *echo.Group, config config.Config) {
	// Public routes
	g.GET("/company/:id", WithContext(handler.GetCompany))
}

func (r *CompanyRouteProvider) AddPrivateRoutes(g *echo.Group, config config.Config) {
	// Private routes (requiring authentication)
}
