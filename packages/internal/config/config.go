package config

import (
	"fmt"
	"os"
)

type PGConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (c PGConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode,
	)
}

func GetPGConfig() PGConfig {
	return PGConfig{
		Host:     getenv("PG_HOST", "localhost"),
		Port:     getenv("PG_PORT", "5432"),
		User:     getenv("PG_USER", "gobasic"),
		Password: getenv("PG_PASSWORD", "gobasic"),
		Database: getenv("PG_DATABASE", "gobasic"),
		SSLMode:  getenv("PG_SSLMODE", "disable"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

