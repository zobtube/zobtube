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
	Metadata struct {
		Type              string // "filesystem" or "s3"
		Path              string // filesystem root path
		S3Bucket          string
		S3Region          string
		S3Prefix          string
		S3Endpoint        string
		S3AccessKeyID     string
		S3SecretAccessKey string
	}
	// DefaultLibraryID is set after bootstrap; used as the default upload target for videos.
	DefaultLibraryID string
	Authentication   bool
}

func New(logger *zerolog.Logger, serverBind, dbDriver, dbConnstring, mediaPath string, meta MetadataParams) (*Config, error) {
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

	if err := applyMetadataDefaults(cfg, meta); err != nil {
		return cfg, err
	}

	logger.Info().
		Str("db-driver", cfg.DB.Driver).
		Str("server-bind", cfg.Server.Bind).
		Str("media-path", cfg.Media.Path).
		Str("metadata-type", cfg.Metadata.Type).
		Str("metadata-path", cfg.MetadataPathDisplay()).
		Msg("valid configuration found")

	return cfg, nil
}
