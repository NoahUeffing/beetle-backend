package router

import (
	"beetle/internal/config"
	"beetle/internal/handler"
	"beetle/internal/router/middleware"

	"github.com/labstack/echo/v4"
)

type ProductRouteProvider struct{}

func (r *ProductRouteProvider) AddPublicRoutes(g *echo.Group, config config.Config) {
	// Public routes
	g.GET("/product/license/:id", WithContext(handler.GetProductLicense))
	g.GET("/product/licenses", WithContext(handler.GetLicenses, middleware.Paginate))
	g.GET("/product/dosage-forms", WithContext(handler.GetDosageForms, middleware.Paginate))
	g.GET("/product/dosage-form/:id", WithContext(handler.GetDosageFormByID))
	g.GET("/product/license/submission-types", WithContext(handler.GetSubmissionTypes, middleware.Paginate))
	g.GET("/product/license/submission-type/:id", WithContext(handler.GetSubmissionTypeByID))
}

func (r *ProductRouteProvider) AddPrivateRoutes(g *echo.Group, config config.Config) {
	// Private routes (requiring authentication)
}
