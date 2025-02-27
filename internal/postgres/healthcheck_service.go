package postgres

import (
	"beetle/internal/healthcheck"

	"gorm.io/gorm"
)

type HealthCheckService struct {
	DB *gorm.DB
}

func (s *HealthCheckService) Check() healthcheck.Status {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return healthcheck.Status{
			Name: "database",
			Up:   false,
		}
	}

	if err := sqlDB.Ping(); err != nil {
		return healthcheck.Status{
			Name: "database",
			Up:   false,
		}
	}

	return healthcheck.Status{
		Name: "database",
		Up:   true,
	}
}
