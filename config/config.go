package config

import (
	"github.com/caarlos0/env"
)

// Config - Struct for loading/storing service cofiguration information
type Config struct {
	ServerName     string `env:"SERVER_NAME" envDefault:"application.metadata.exercise"`
	HTTPListenHost string `env:"SERVER_HOST" envDefault:"127.0.0.1"`
	HTTPListenPort string `env:"SERVER_PORT" envDefault:"8071"`
}

// GetConfig - Method for getting config and laoding from envrionment variables
func GetConfig() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return cfg, err
}
