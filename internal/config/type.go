package config

import (
	"errors"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

var ErrNoDbDriverSet = errors.New("ZT_DB_DRIVER is not set")
var ErrNoDbConnStringSet = errors.New("ZT_DB_CONNSTRING is not set")
var ErrNoMediaPathSet = errors.New("ZT_MEDIA_PATH is not set")

type Config struct {
	Server struct {
		Bind string `yaml:"bind" envconfig:"ZT_SERVER_BIND"`
	}
	DB struct {
		Driver     string `yaml:"driver" envconfig:"ZT_DB_DRIVER"`
		Connstring string `yaml:"connstring" envconfig:"ZT_DB_CONNSTRING"`
	} `yaml:"db"`
	Media struct {
		Path string `yaml:"path" envconfig:"ZT_MEDIA_PATH"`
	} `yaml:"media"`
}

func New(configPath string) (*Config, error) {
	cfg := &Config{}

	if _, err := os.Stat(configPath); err == nil {
		f, err := os.Open(configPath)
		if err != nil {
			return cfg, err
		}
		defer f.Close()

		decoder := yaml.NewDecoder(f)
		err = decoder.Decode(cfg)
		if err != nil {
			return cfg, err
		}
	}

	err := envconfig.Process("zt", cfg)
	if err != nil {
		return cfg, err
	}

	// pre flight checks
	if cfg.DB.Driver == "" {
		return cfg, ErrNoDbDriverSet
	}

	if cfg.DB.Connstring == "" {
		return cfg, ErrNoDbConnStringSet
	}

	if cfg.Media.Path == "" {
		return cfg, ErrNoMediaPathSet
	}

	if cfg.Server.Bind == "" {
		cfg.Server.Bind = "127.0.0.1:8080"
	}

	return cfg, nil
}
