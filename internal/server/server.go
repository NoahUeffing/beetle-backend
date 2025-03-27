package server

import (
	_ "beetle/docs" // Required for Swagger docs
	"beetle/internal/config"
	"beetle/internal/handler"
	"beetle/internal/router"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Echo   *echo.Echo
	Config config.Config
}

// @title BEETLE API
// @version 1.0

// @BasePath /v1
// @securityDefinitions.apikey JWTToken
// @in header
// @name Authorization

// TODO: Above security definition is not working, auth does not work with curl either

func New(config *config.Config, contextConfig handler.ContextConfig) *Server {
	s := &Server{
		Echo:   echo.New(),
		Config: *config,
	}

	s.Echo.HideBanner = true
	s.Echo.HidePort = config.Logs.HideStartupMessage
	s.addLogging()

	s.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:*", "*.local:*", "*.localhost:*"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowMethods:     []string{http.MethodPost, http.MethodPut, http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodGet, http.MethodPatch},
		AllowCredentials: true,
	}))

	s.Echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "/swagger") // see https://github.com/swaggo/echo-swagger#example
		},
	}))

	s.setContext(contextConfig)
	s.setErrorHandler()

	r := router.New(s.Config)
	r.AddRoutes(
		s.Echo,
		contextConfig.AuthService.GetMiddleware(),
	)

	return s
}

func (s *Server) Start() {
	s.Echo.Logger.Fatal(s.Echo.Start(":8080"))
}

func (s *Server) addLogging() {
	loggerConfig := middleware.LoggerConfig{}
	// For available template options see https://echo.labstack.com/middleware/logger/#configuration
	if s.Config.Logs.JSON {
		// Default message for now, but wanted to be explicit so that future changes are easier to see.
		loggerConfig.Format = `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n"
	} else {
		// Close to W3C format, primarily to be used for local dev
		loggerConfig.Format = `${time_rfc3339} ${remote_ip} ${method} ${uri} ${status} ${latency_human} Error: ${error}` + "\n"
	}
	s.Echo.Use(middleware.LoggerWithConfig(loggerConfig))
	s.Echo.Logger.SetLevel(s.Config.Logs.Level)
}

func (s *Server) setContext(contextConfig handler.ContextConfig) {
	s.Echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) error {
			return next(&handler.Context{
				Context:       echoContext,
				ContextConfig: contextConfig,
			})
		}
	})
}

func (s *Server) setErrorHandler() {
	s.Echo.HTTPErrorHandler = func(err error, c echo.Context) {
		s.Echo.DefaultHTTPErrorHandler(handler.WrapError(err), c)
	}
}
