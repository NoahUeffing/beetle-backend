package config

import (
	"beetle/internal/env"
	"log"

	"github.com/spf13/viper"
)

var configFolders = []string{
	"../configs",
	"./configs",
	"./",
}

func (c *Config) loadFromFile() error {
	viper.SetConfigName(c.Env)
	viper.SetConfigType("yaml")
	for _, dir := range configFolders {
		viper.AddConfigPath(dir)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Error loading config file:\n" + err.Error())
		return err
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatalln("Error un-marshalling config file:\n" + err.Error())
		return err
	}

	return nil
}

func (c *Config) LoadFromEnv() error {
	c.DB.Read = env.Get("BEETLE_DB_READ", c.DB.Read)
	c.DB.Write = env.Get("BEETLE_DB_WRITE", c.DB.Write)
	c.Auth.Secret = env.Get("JWT_SECRET", c.Auth.Secret)
	c.MigrationDir = env.Get("MIGRATION_DIR", c.MigrationDir)
	c.Email.ApiKey = env.Get("MAILERSEND_TOKEN", c.Email.ApiKey)
	c.Email.From = env.Get("MAILERSEND_USER", c.Email.From)
	return nil
}

func Load() *Config {
	c := &Config{}
	c.Env = env.Get("ENV", "dev")
	c.loadFromFile()
	c.LoadFromEnv()
	return c
}
