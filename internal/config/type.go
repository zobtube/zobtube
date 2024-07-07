package config

import (
	"errors"
	"os"
)

type Config struct {
	DbDriver     string
	DbConnstring string
	MediaFolder  string
}

func New() (*Config, error) {
	cfg := &Config{}

	cfg.DbDriver = os.Getenv("ZT_DB_DRIVER")
	if cfg.DbDriver == "" {
		return cfg, errors.New("ZT_DB_DRIVER is not set")
	}

	cfg.DbConnstring = os.Getenv("ZT_DB_CONNSTRING")
	if cfg.DbConnstring == "" {
		return cfg, errors.New("ZT_DB_CONNSTRING is not set")
	}

	cfg.MediaFolder = os.Getenv("ZT_MEDIA")
	if cfg.MediaFolder == "" {
		return cfg, errors.New("ZT_MEDIA is not set")
	}

	return cfg, nil
}
