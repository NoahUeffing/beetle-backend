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
// @description This is the Beetle API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer: ` prefix, e.g. "Bearer abcde12345".

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
