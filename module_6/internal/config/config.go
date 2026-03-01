package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env                string `env:"ENV" envDefault:"development"`
	MongoConnectionURI string `env:"MONGO_CONNECTION_URI" envDefault:"mongodb://root:root@mongo:27017/?authSource=admin"`
	MongoDbName        string `env:"MONGO_DATABASE_NAME" envDefault:"app"`
	LogLevel           string `env:"LOG_LEVEL" envDefault:"info"`
	TokenTTL           string `env:"TOKEN_TTL" envDefault:"1h"`
	TokenSecret        string `env:"TOKEN_SECRET" envDefault:"secret"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	if cfg.MongoConnectionURI == "" || cfg.MongoDbName == "" {
		return nil, fmt.Errorf("missing MongoDB configuration")
	}
	return cfg, nil
}
