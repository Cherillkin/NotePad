package config

import (
	"github.com/caarlos0/env"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DBPORT        string `env:"DB_PORT"`
	DBHost        string `env:"DB_HOST"`
	DBName        string `env:"DB_NAME"`
	DBPassword    string `env:"POSTGRES_PASSWORD"`
	DBUser        string `env:"DB_USER"`
	DBSSLMode     string `env:"DB_SSLMODE"`
	RedisAddr     string `env:"REDIS_ADDR"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"`
}

func NewEnvConfig() *EnvConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load the .env file: %e", err)
	}

	config := &EnvConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("Unable to load variables from the .env: %e", err)
	}

	return config
}
