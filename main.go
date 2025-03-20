package main

import (
	"log"

	_ "beetle/docs" // Required for Swagger docs
	"beetle/internal/auth"
	"beetle/internal/config"
	"beetle/internal/handler"
	"beetle/internal/healthcheck"
	"beetle/internal/postgres"
	"beetle/internal/server"
	"beetle/internal/validation"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"gorm.io/gorm"
)

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

	userService := &postgres.UserService{
		ReadDB:      gormReadDB,
		WriteDB:     gormWriteDB,
		AuthService: authService,
	}

	cc := handler.ContextConfig{
		AuthService: authService,

		UserService: userService,

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

func getSqlxDBFromGorm(gormDB *gorm.DB) (*sqlx.DB, error) {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}
	return sqlx.NewDb(sqlDB, "postgres"), nil
}

func migrate(dir string, db *sqlx.DB) {
	if err := goose.Run("up", db.DB, dir); err != nil {
		log.Fatalf("Unable to complete migrations, error:\n%+v\n", err)
	}
}
