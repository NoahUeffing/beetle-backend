package router

import (
	"beetle/internal/config"
	"beetle/internal/handler"
	"beetle/internal/router/middleware"

	"github.com/labstack/echo/v4"
)

type productRouteProvider struct{}

func (r *productRouteProvider) AddPublicRoutes(g *echo.Group, config config.Config) {
	// Public routes
	g.GET("/product/license/:id", WithContext(handler.GetProductLicense))
	g.GET("/product/dosage-forms", WithContext(handler.GetDosageForms, middleware.Paginate))
	g.GET("/product/license/submission-types", WithContext(handler.GetSubmissionTypes, middleware.Paginate))
}

func (r *productRouteProvider) AddPrivateRoutes(g *echo.Group, config config.Config) {
	// Private routes (requiring authentication)
}
