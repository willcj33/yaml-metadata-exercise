package config

import (
	"github.com/caarlos0/env"
)

// Config - Struct for loading/storing service cofiguration information
type Config struct {
	ServerName       string   `env:"SERVER_NAME" envDefault:"application.metadata.exercise"`
	HTTPListenHost   string   `env:"SERVER_HOST" envDefault:"127.0.0.1"`
	HTTPListenPort   string   `env:"SERVER_PORT" envDefault:"8071"`
	StorageMode      string   `env:"STORAGE_MODE" envDefault:"multiple"`
	IndexName        string   `env:"INDEX_NAME" envDefault:"applicationMetadata.bleve"`
	IdentifierFields []string `env:"IDENTIFIER_FIELDS" envSeparator:"," envDefault:"title,source"`
}

// GetConfig - Method for getting config and laoding from envrionment variables
func GetConfig() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return cfg, err
}
