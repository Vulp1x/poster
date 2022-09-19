package config

import (
	"context"
	"io"
	"os"

	"github.com/inst-api/poster/internal/postgres"
	"github.com/inst-api/poster/internal/sessions"
	"github.com/inst-api/poster/pkg/logger"
	"gopkg.in/yaml.v3"
)

const (
	localConfigFilePath = "deploy/configs/values_local.yaml"
	prodConfigFilePath  = "deploy/configs/values_production.yaml"
)

// Config represents application configuration.
type Config struct {
	Listen   ServerConfig
	Logger   logger.Configuration   `yaml:"logger"`
	Postgres postgres.Configuration `yaml:"postgres"`
	Security sessions.Configuration `yaml:"session"`
}

// ServerConfig represents configuration of server location
type ServerConfig struct {
	BindIP string `yaml:"bind_ip"`
	Port   string `yaml:"port"`
}

// ParseConfiguration parses configuration from goadmin_config.yml.
func (c *Config) ParseConfiguration(local bool) error {
	c.Default()

	configFilePath := prodConfigFilePath
	if local {
		configFilePath = localConfigFilePath
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		logger.Errorf(context.Background(), "failed to open config file at %s: %v", configFilePath, err)
		return nil
		// return fmt.Errorf("failed to open config file %s: %v", configFilePath, err)
	}

	data, _ := io.ReadAll(configFile)

	logger.Infof(context.Background(), "starting with config from %s", configFilePath)

	return yaml.Unmarshal(data, c)
}

// Default sets default values in config variables.
func (c *Config) Default() {
	c.Listen = ServerConfig{BindIP: "0.0.0.0", Port: "8090"}
	c.Logger.Default()
	c.Postgres.Default()
	c.Security.Default()
}
