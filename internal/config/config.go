package config

import (
	"github.com/labstack/gommon/log"
)

type Config struct {
	Env          string
	MigrationDir string
	Logs         LogsConfig
	DB           DBsConfig
	Auth         AuthConfig
	Email        EmailConfig
}

type DBsConfig struct {
	Write string // DB Connection String
	Read  string // DB Connection String
}

type LogsConfig struct {
	Level              log.Lvl
	JSON               bool
	HideStartupMessage bool
}

type AuthConfig struct {
	Secret string
}

type EmailConfig struct {
	ApiKey string
	From   string
}
