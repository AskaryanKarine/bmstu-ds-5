package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppEnv      string `env:"APP_ENV" env-default:"test"`
	PostgresDSN string `env:"POSTGRES_DSN"`
	Port        int    `env:"PORT" env-default:"8080"`
}

func ReadConfig(localPath string) (Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, err
	}

	if cfg.AppEnv != "test" {
		return cfg, nil
	}

	err = cleanenv.ReadConfig(localPath, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
