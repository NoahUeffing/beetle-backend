package main

import (
	"log"

	"beetle/internal/auth"
	"beetle/internal/config"
	"beetle/internal/handler"
	"beetle/internal/healthcheck"
	"beetle/internal/postgres"
	"beetle/internal/server"
	"beetle/internal/validation"
	_ "beetle/swaggergenerated" // Required for Swagger docs

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
)

// @title Beetle API
// @version 1.0

// @BasePath /v1
// @securityDefinitions.apikey JWTToken
// @in header
// @name Authorization

// todo: integrate api  https://health-products.canada.ca/api/documentation/lnhpd-documentation-en.html

func main() {
	config := config.Load()
	readDB, writeDB := postgres.Open(config.DB)
	defer postgres.Close(readDB)
	defer postgres.Close(writeDB)
	gormReadDB, gormWriteDB := postgres.OpenGorm(config.DB)

	// TODO: Determine how to migrate down when downgrading/rolling back
	migrate(config.MigrationDir, writeDB)

	authService := auth.New(config.Auth)

	cc := handler.ContextConfig{
		AuthService: authService,
		UserService: &postgres.UserService{
			ReadDB:      gormReadDB,
			WriteDB:     gormWriteDB,
			AuthService: authService,
		},
		ProductService: &postgres.ProductService{
			ReadDB:  gormReadDB,
			WriteDB: gormWriteDB,
		},
		ValidationService: *validation.New(),
		HealthCheckServices: []healthcheck.IHealthCheckService{
			&postgres.HealthCheckService{
				Read:  gormReadDB,
				Write: gormWriteDB,
			},
		},
	}

	s := server.New(config, cc)
	s.Start()
}

func migrate(dir string, db *sqlx.DB) {
	if err := goose.Run("up", db.DB, dir); err != nil {
		log.Fatalf("Unable to complete migrations, error:\n%+v\n", err)
	}
}
