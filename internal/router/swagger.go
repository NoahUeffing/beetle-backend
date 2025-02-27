package router

import (
	"beetle/internal/config"
	"fmt"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func AddSwaggerRoutes(s *echo.Echo, config config.Config) {
	s.GET("/swagger/*", echoSwagger.WrapHandler)
	if !config.Logs.HideStartupMessage {
		fmt.Println("Swagger UI available on port http://localhost:8080/swagger/index.html")
	}
}
