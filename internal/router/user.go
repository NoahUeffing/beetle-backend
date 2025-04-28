package router

import (
	"beetle/internal/config"
	"beetle/internal/handler"
	"beetle/internal/router/middleware"

	"github.com/labstack/echo/v4"
)

type UserRouteProvider struct{}

func (r *UserRouteProvider) AddPublicRoutes(g *echo.Group, config config.Config) {
	// Public user routes
	g.POST("/user", WithContext(handler.CreateUser))
}

func (r *UserRouteProvider) AddPrivateRoutes(g *echo.Group, config config.Config) {
	// Private user routes (requiring authentication)
	g.GET("/user/:id", WithContext(handler.GetUser, middleware.Auth))
	g.PUT("/user/:id", WithContext(handler.UpdateUser, middleware.Auth))
}
