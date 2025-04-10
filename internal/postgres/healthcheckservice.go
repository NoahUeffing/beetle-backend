package postgres

import (
	"beetle/internal/healthcheck"

	"gorm.io/gorm"
)

// UserService represents a PostgreSQL implementation of myapp.UserService.
type HealthCheckService struct {
	Read  *gorm.DB
	Write *gorm.DB
}

const HEALTH_CHECK_NAME = "Postgres Database Connection"

// User returns a user for a given id.
func (hcs *HealthCheckService) Check() healthcheck.Status {
	status := healthcheck.Status{
		Name:     HEALTH_CHECK_NAME,
		Up:       true,
		Messages: []string{},
	}

	sqlDB, err := hcs.Read.DB()
	if err != nil {
		status.Messages = append(status.Messages, "Read DB Unavailable: "+err.Error())
		status.Up = false
	} else if err := sqlDB.Ping(); err != nil {
		status.Messages = append(status.Messages, "Read DB Unavailable: "+err.Error())
		status.Up = false
	}

	sqlDB, err = hcs.Write.DB()
	if err != nil {
		status.Messages = append(status.Messages, "Write DB Unavailable: "+err.Error())
		status.Up = false
	} else if err := sqlDB.Ping(); err != nil {
		status.Messages = append(status.Messages, "Write DB Unavailable: "+err.Error())
		status.Up = false
	}

	return status
}
