package config

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/inst-api/poster/internal/postgres"
	"github.com/inst-api/poster/internal/sessions"
	"github.com/inst-api/poster/pkg/logger"
	"gopkg.in/yaml.v3"
)

const (
	localConfigFilePath       = "deploy/configs/values_local.yaml"
	localDockerConfigFilePath = "deploy/configs/values_docker.yaml"
	prodConfigFilePath        = "deploy/configs/values_production.yaml"

	localConfigMode       = "local"
	localDockerConfigMode = "local_docker"
	prodConfigMode        = "prod"
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
	BindIP   string `yaml:"bind_ip"`
	Port     string `yaml:"port"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

// ParseConfiguration parses configuration from values_*.yaml
func (c *Config) ParseConfiguration(configMode string) error {
	c.Default()

	var configFilePath string
	switch {
	case configMode == localConfigMode:
		configFilePath = localConfigFilePath
	case configMode == localDockerConfigMode:
		configFilePath = localDockerConfigFilePath
	case configMode == prodConfigMode:
		configFilePath = prodConfigFilePath
	default:
		return fmt.Errorf(
			"unexpected config mode: '%s', expected one of ['%s', '%s', '%s']",
			configMode,
			localConfigMode,
			localDockerConfigMode,
			prodConfigMode,
		)
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
