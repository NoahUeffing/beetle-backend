package router

import (
	"beetle/internal/config"
	"beetle/internal/handler"

	"github.com/labstack/echo/v4"
)

type productRouteProvider struct{}

func (r *productRouteProvider) AddPublicRoutes(g *echo.Group, config config.Config) {
	// Public routes
	g.GET("/product/license/:id", WithContext(handler.GetProductLicense))
}

func (r *productRouteProvider) AddPrivateRoutes(g *echo.Group, config config.Config) {
	// Private routes (requiring authentication)
}
