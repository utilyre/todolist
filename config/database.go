package config

import (
	"os"

	"go.uber.org/fx"
)

type DatabaseConfig struct {
	Path string
}

func NewDatabaseConfig(lc fx.Lifecycle) DatabaseConfig {
	return DatabaseConfig{
		Path: os.Getenv("DB_PATH"),
	}
}
