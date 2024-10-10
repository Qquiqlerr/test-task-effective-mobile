package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Logger   LoggerConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL"`
}

type ServerConfig struct {
	Port string `env:"PORT"`
}

// MustLoad загружает конфигурацию из файла .env или выдаёт панику
func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("failed to parse .env file. Error: %e", err))
	}
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(fmt.Sprintf("failed to parse .env file. Error: %e", err))
	}
	return &cfg
}
