package router

import (
	"beetle/internal/config"
	_ "beetle/swaggergenerated" // Generated by `swag init`

	"github.com/labstack/echo/v4"
)

type RouteProvider interface {
	AddPublicRoutes(g *echo.Group, config config.Config)
	AddPrivateRoutes(g *echo.Group, config config.Config)
}

type Router struct {
	RouteProviders []RouteProvider
	Config         config.Config
}

func New(config config.Config) *Router {
	return &Router{
		RouteProviders: []RouteProvider{
			&TokenRouteProvider{},
			&HealthcheckRouteProvider{},
			&UserRouteProvider{},
			&CompanyRouteProvider{},
			&ProductRouteProvider{},
		},
		Config: config,
	}
}

// AddRoutes adds all routes to the Echo instance
func (r *Router) AddRoutes(server *echo.Echo, authMiddlewares ...echo.MiddlewareFunc) {
	v1public := server.Group("/v1")
	v1private := server.Group("/v1", authMiddlewares...)

	for _, rp := range r.RouteProviders {
		rp.AddPrivateRoutes(v1private, r.Config)
	}

	for _, rp := range r.RouteProviders {
		rp.AddPublicRoutes(v1public, r.Config)
	}

	AddSwaggerRoutes(server, r.Config)
}
