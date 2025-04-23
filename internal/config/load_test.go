package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromEnv(t *testing.T) {
	// Setup test environment variables
	os.Setenv("BEETLE_DB_READ", "test-read-db")
	os.Setenv("BEETLE_DB_WRITE", "test-write-db")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("MIGRATION_DIR", "test-migrations")
	defer func() {
		os.Unsetenv("BEETLE_DB_READ")
		os.Unsetenv("BEETLE_DB_WRITE")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("MIGRATION_DIR")
	}()

	config := &Config{}
	err := config.loadFromEnv()
	require.NoError(t, err)

	assert.Equal(t, "test-read-db", config.DB.Read)
	assert.Equal(t, "test-write-db", config.DB.Write)
	assert.Equal(t, "test-secret", config.Auth.Secret)
	assert.Equal(t, "test-migrations", config.MigrationDir)
}
