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

func NewUserRouteProvider(taskHandler *handler.UserHandler) *userRouteProvider {
	return &userRouteProvider{
		userHandler: taskHandler,
	}
}

func (r *userRouteProvider) AddPublicRoutes(g *echo.Group, config config.Config) {
	// We need to adapt the handler methods to work with Echo
	// This will be implemented later
}

func (r *userRouteProvider) AddPrivateRoutes(g *echo.Group, config config.Config) {
	// No private routes for now
}
