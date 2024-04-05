package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP   HTTPConfig
	PG     PGConfig
	Logger LoggerConfig
}

type (
	HTTPConfig struct {
		Host    string `env:"HTTP_HOST" yaml:"host"`
		Port    string `env:"HTTP_PORT" yaml:"port"`
		Address string
	}

	PGConfig struct {
		URL      string
		User     string `env:"POSTGRES_USER"     yaml:"user"`
		Password string `env:"POSTGRES_PASSWORD" yaml:"password"`
		Host     string `env:"POSTGRES_HOST"     yaml:"host"`
		Port     uint16 `env:"POSTGRES_PORT"     yaml:"port"`
		DB       string `env:"POSTGRES_DB"       yaml:"DB"`
		SSLMode  string `env:"POSTGRES_SSL_MODE" yaml:"SSLMode"`
	}

	LoggerConfig struct {
		LogFilePath string `env:"LOG_FILE_PATH" yaml:"logFilePath"`
		Level       string `env:"LOG_LVL"       yaml:"level"`
	}
)

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from environment variables: %w", err)
	}
	cfg.PG.URL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.PG.User, cfg.PG.Password, cfg.PG.Host, cfg.PG.Port, cfg.PG.DB, cfg.PG.SSLMode)

	return cfg, nil
}
