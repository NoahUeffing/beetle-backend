package config

import (
	"github.com/labstack/gommon/log"
)

type Config struct {
	Env          string
	MigrationDir string
	Logs         LogsConfig
	DB           DBConfig
	Auth         AuthConfig
}

type LogsConfig struct {
	Level              log.Lvl
	JSON               bool
	HideStartupMessage bool
}

type AuthConfig struct {
	Secret string
}
