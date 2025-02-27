package server

import (
	"fmt"
	"log"
	"os"

	"beetle/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title Beetle API
// @version 1.0
// @description This is the REST API for the Beetle platform
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email your-email@domain.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

type Server struct {
	e  *echo.Echo
	cc handler.ContextConfig
}

func (s *Server) contextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := &handler.Context{
			Context:       c,
			ContextConfig: s.cc,
		}
		return next(ctx)
	}
}

func (s *Server) wrap(h func(handler.Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		if ctx, ok := c.(*handler.Context); ok {
			return h(*ctx)
		}
		return echo.ErrInternalServerError
	}
}

func New(cc handler.ContextConfig) *Server {
	s := &Server{
		e:  echo.New(),
		cc: cc,
	}

	// Middleware
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())
	s.e.Use(middleware.CORS())
	s.e.Use(s.contextMiddleware)

	// Routes
	s.e.GET("/healthcheck", s.wrap(handler.HealthCheck))

	// API v1 routes
	v1 := s.e.Group("/v1")
	v1.POST("/tokens", s.wrap(handler.AuthTokenCreate))

	return s
}

func (s *Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Starting server on %s", addr)
	return s.e.Start(addr)
}
