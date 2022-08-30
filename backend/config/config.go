package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Env        string `env:"BASEAPP_ENV" envDefault:"dev"`
	Port       int    `env:"PORT" envDefault:"80"`
	DBHost     string `env:"DB_HOST" envDefault:"127.0.0.1"`
	DBPort     int    `env:"DB_PORT" envDefault:"33306"`
	DBUser     string `env:"DB_USER" envDefault:"user"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"password"`
	DBName     string `env:"DB_NAME" envDefault:"baseapp"`
}

func New() (*Config, error) {
	cfg := &Config{}
	// https://github.com/caarlos0/env#example
	// os.Getenv()との違いは，タグによってデフォルト値を設定できること
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
