package postgres

import (
	"bytes"
	"fmt"
	"time"
)

// Configuration represents database configuration.
type Configuration struct {
	MigrationsDir string        `yaml:"migrations dir"`
	Timeout       time.Duration `yaml:"timeout"`
	MaxConns      int32         `yaml:"maxConns"`
	Host          string        `yaml:"host"`
	Port          int           `yaml:"port"`
	User          string        `yaml:"user"`
	Password      string        `yaml:"password"`
	Database      string        `yaml:"database"`
}

// Default sets default values in config variables.
func (c *Configuration) Default() {
	c.Host = "postgres"
	c.Port = 5432
	c.Database = "insta_poster"
	c.User = "docker"
	c.Password = "docker"
	c.MigrationsDir = "migrations"
	c.Timeout = 5 * time.Second
	c.MaxConns = 25
}

func (c *Configuration) buildFullDsn() string {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("user=%s password=%s database=%s host=%s port=%d ",
		c.User, c.Password, c.Database, c.Host, c.Port))

	if c.Timeout > 0 {
		b.WriteString(fmt.Sprintf("connect_timeout=%d ", c.Timeout.Milliseconds()/1000))
	}

	if c.MaxConns > 0 {
		b.WriteString(fmt.Sprintf("pool_max_conns=%d ", c.MaxConns))
	}

	return b.String()
}
