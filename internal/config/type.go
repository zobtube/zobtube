package config

import (
	"errors"

	"github.com/rs/zerolog"
)

type Config struct {
	Server struct {
		Bind string
	}
	DB struct {
		Driver     string
		Connstring string
	}
	Media struct {
		Path string
	}
	Authentication bool
}

func New(logger *zerolog.Logger, serverBind, dbDriver, dbConnstring, mediaPath string) (*Config, error) {
	cfg := &Config{}
	cfg.Server.Bind = serverBind
	cfg.DB.Driver = dbDriver
	cfg.DB.Connstring = dbConnstring
	cfg.Media.Path = mediaPath

	// pre flight checks
	if cfg.DB.Driver == "" {
		return cfg, errors.New("ZT_DB_DRIVER is not set")
	}

	if cfg.DB.Connstring == "" {
		return cfg, errors.New("ZT_DB_CONNSTRING is not set")
	}

	if cfg.Media.Path == "" {
		return cfg, errors.New("ZT_MEDIA_PATH is not set")
	}

	logger.Info().
		Str("db-driver", cfg.DB.Driver).
		Str("server-bind", cfg.Server.Bind).
		Str("media-path", cfg.Media.Path).
		Msg("valid configuration found")

	return cfg, nil
}
