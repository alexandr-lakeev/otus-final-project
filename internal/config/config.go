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
		Server    ServerConf    `config:"server"`
		Previewer PreviewerConf `config:"previewer"`
		Logger    LoggerConf    `config:"logger"`
	}

	ServerConf struct {
		BindAddress  string        `yaml:"http_bind_address" config:"http_bind_address"`
		ReadTimeout  time.Duration `yaml:"http_read_timeout" config:"http_read_timeout"`
		WriteTimeout time.Duration `yaml:"http_write_timeout" config:"http_write_timeout"`
		IdleTimeout  time.Duration `yaml:"http_idle_timeout" config:"http_idle_timeout"`
	}

	PreviewerConf struct {
		RequestTimeout time.Duration `yaml:"request_timeout" config:"request_timeout"`
		CacheSize      int           `yaml:"cache_size" config:"cache_size"`
		CacheDir       string        `yaml:"cache_dir" config:"cache_dir"`
	}

	LoggerConf struct {
		Env   string `config:"ENV"`
		Level string `yaml:"level"  config:"level"`
	}
)

func NewConfig(configFile string) (*Config, error) {
	cfg := Config{
		Server: ServerConf{
			BindAddress: ":8080",
		},
	}

	if err := confita.NewLoader(
		env.NewBackend(),
		file.NewBackend(configFile),
	).Load(context.Background(), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
