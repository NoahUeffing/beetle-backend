package main

import (
	"fmt"
	"log"

	_ "beetle/docs" // Required for Swagger docs
	"beetle/internal/auth"
	"beetle/internal/config"
	"beetle/internal/handler"
	"beetle/internal/postgres"
	"beetle/internal/router"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Beetle API
// @version 1.0

// @BasePath /v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Enter the token with the `Bearer: ` prefix, e.g. "Bearer abcde12345".

// TODO: Above security definition is not working
func main() {
	// Initialize database connection
	dbConfig := config.NewDBConfig()
	fmt.Println("Host")
	fmt.Println(dbConfig.Host)
	db, err := config.NewDatabase(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize auth service
	authService := auth.New(config.AuthConfig{
		Secret: "your-secret-key", // TODO: Load from environment/config
	})

	// Initialize user service
	userService := &postgres.UserService{
		ReadDB:      db,
		WriteDB:     db,
		AuthService: authService,
	}

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService)

	// Initialize router
	appConfig := config.Config{} // Create a proper config if needed
	r := router.New(appConfig, userHandler, authHandler)

	// Add routes to Echo
	r.AddRoutes(e)

	// Add Swagger handler
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	log.Printf("Starting server on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
