package postgres

import (
	"beetle/internal/healthcheck"

	"gorm.io/gorm"
)

type HealthCheckService struct {
	Read  *gorm.DB
	Write *gorm.DB
}

func (s *HealthCheckService) Check() healthcheck.Status {
	sqlDB, err := s.Read.DB()
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
