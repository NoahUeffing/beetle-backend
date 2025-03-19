package main

import (
	"fmt"
	"log"

	"beetle/internal/auth"
	"beetle/internal/config"
	"beetle/internal/handler"
	"beetle/internal/postgres"
	"beetle/internal/server"
	"beetle/internal/validation"
)

// @title Beetle API
// @version 1.0
// @description This is the REST API for the Beetle platform
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@beetle.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1
// @schemes http https

// @securityDefinitions.apikey JWTToken
// @in header
// @name Authorization
// @description Bearer token authentication

func main() {
	// Initialize database connection
	dbConfig := config.NewDBConfig()
	fmt.Println("Host")
	fmt.Println(dbConfig.Host)
	db, err := config.NewDatabase(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

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

	// Initialize validation service
	validationService := validation.New()

	// Create context configuration
	contextConfig := handler.ContextConfig{
		AuthService:       authService,
		UserService:       userService,
		ValidationService: *validationService,
	}

	// Initialize server
	s := server.New(contextConfig)

	// Start server
	log.Printf("Starting server on :8080")
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
