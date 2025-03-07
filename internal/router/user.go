package router

import (
	"beetle/internal/config"
	"beetle/internal/handler"

	"github.com/labstack/echo/v4"
)

type userRouteProvider struct {
	userHandler *handler.UserHandler
}

// TODO:

func NewUserRouteProvider(userHandler *handler.UserHandler) *userRouteProvider {
	return &userRouteProvider{
		userHandler: userHandler,
	}
}

func (r *userRouteProvider) AddPublicRoutes(g *echo.Group, config config.Config) {
	// Public user routes
	g.POST("/user", r.userHandler.CreateUser)
}

func (r *userRouteProvider) AddPrivateRoutes(g *echo.Group, config config.Config) {
	// Private user routes (requiring authentication)
	// None for now
}
