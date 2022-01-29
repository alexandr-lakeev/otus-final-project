package config

import (
	"context"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

type (
	Config struct {
		Server    ServerConf    `toml:"server"`
		Previewer PreviewerConf `toml:"previewer"`
		Logger    LoggerConf    `toml:"logger"`
	}

	ServerConf struct {
		BindAddress  string        `config:"http_bind_address,require"`
		ReadTimeout  time.Duration `config:"http_read_timeout"`
		WriteTimeout time.Duration `config:"http_write_timeout"`
		IdleTimeout  time.Duration `config:"http_idle_timeout"`
	}

	PreviewerConf struct {
		CacheSize int `config:"cache_size"`
	}

	LoggerConf struct {
		Env   string `config:"ENV"`
		Level string `config:"level"`
	}
)

func NewConfig(configFile string) (*Config, error) {
	cfg := Config{
		Server: ServerConf{
			BindAddress: ":8080",
		},
	}

	err := confita.NewLoader(
		env.NewBackend(),
		file.NewBackend(configFile),
	).Load(context.Background(), &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
