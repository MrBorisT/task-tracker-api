package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
	Port       string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		Port:       os.Getenv("APP_PORT"),
	}

	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" ||
		cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBSSLMode == "" {
		return nil, fmt.Errorf("missing required database configuration environment variables")
	}

	if cfg.Port == "" {
		cfg.Port = ":8080"
	} else if !strings.HasPrefix(cfg.Port, ":") {
		cfg.Port = ":" + cfg.Port
	}

	return cfg, nil
}
