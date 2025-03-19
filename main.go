package main

import (
	"fmt"
	"log"

	"beetle/internal/auth"
	"beetle/internal/config"
	"beetle/internal/handler"
	"beetle/internal/postgres"
	"beetle/internal/router"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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

	// Initialize user handler
	userHandler := handler.NewUserHandler(userService)

	// Initialize router
	appConfig := config.Config{} // Create a proper config if needed
	r := router.New(appConfig, userHandler)

	// Add routes to Echo
	r.AddRoutes(e)

	// Start server
	log.Printf("Starting server on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
