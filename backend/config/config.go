package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Env  string `env:"TODO_ENV" envDefault:"dev"`
	Port int    `env:"PORT" envDefault:"80"`
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