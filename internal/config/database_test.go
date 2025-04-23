package config_test

import (
	"beetle/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBConfig_DSN(t *testing.T) {
	config := &config.DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	assert.Equal(t, expected, config.DSN())
}
