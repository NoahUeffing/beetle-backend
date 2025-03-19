package server

import (
	"fmt"
	"log"
	"os"

	_ "beetle/docs" // Required for Swagger docs
	"beetle/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	e  *echo.Echo
	cc handler.ContextConfig
}

// @title Beetle API
// @version 1.0

// @BasePath /v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Enter the token with the `Bearer: ` prefix, e.g. "Bearer abcde12345".

// TODO: Above security definition is not working, auth does not work with curl either

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

	// Swagger documentation
	s.e.GET("/swagger/*", echoSwagger.WrapHandler)

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
