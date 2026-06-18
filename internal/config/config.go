package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppEnv string `env:"APP_ENV" env-default:"local"`

	HttpHost string `env:"HTTP_HOST" env-default:"0.0.0.0"`
	HttpPort string `env:"HTTP_PORT" env-default:"8080"`

	PostgresHost     string `env:"POSTGRES_HOST" env-required:"true"`
	PostgresPort     string `env:"POSTGRES_PORT" env-required:"true"`
	PostgresDB       string `env:"POSTGRES_DB" env-required:"true"`
	PostgresUser     string `env:"POSTGRES_USER" env-required:"true"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-required:"true"`

	Source1 string `env:"SOURCE1" env-required:"true"`
	Source2 string `env:"SOURCE2" env-required:"true"`
	Source3 string `env:"SOURCE3" env-required:"true"`

	ClientsSource string `env:"CLIENTS_SOURCE" env-required:"true"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	return &cfg
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresDB,
	)
}

func (c *Config) Address() string {
	return fmt.Sprintf("%s:%s", c.HttpHost, c.HttpPort)
}