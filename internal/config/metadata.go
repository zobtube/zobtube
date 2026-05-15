package config

import (
	"errors"
	"fmt"
)

// MetadataParams holds metadata storage settings from CLI, env, or YAML.
type MetadataParams struct {
	Type              string
	Path              string
	S3Bucket          string
	S3Region          string
	S3Prefix          string
	S3Endpoint        string
	S3AccessKeyID     string
	S3SecretAccessKey string
}

func applyMetadataDefaults(cfg *Config, meta MetadataParams) error {
	cfg.Metadata.Type = meta.Type
	cfg.Metadata.Path = meta.Path
	cfg.Metadata.S3Bucket = meta.S3Bucket
	cfg.Metadata.S3Region = meta.S3Region
	cfg.Metadata.S3Prefix = meta.S3Prefix
	cfg.Metadata.S3Endpoint = meta.S3Endpoint
	cfg.Metadata.S3AccessKeyID = meta.S3AccessKeyID
	cfg.Metadata.S3SecretAccessKey = meta.S3SecretAccessKey

	if cfg.Metadata.Type == "" {
		cfg.Metadata.Type = "filesystem"
	}

	switch cfg.Metadata.Type {
	case "filesystem":
		if cfg.Metadata.Path == "" {
			cfg.Metadata.Path = cfg.Media.Path
		}
	case "s3":
		if cfg.Metadata.S3Bucket == "" {
			return errors.New("metadata s3 bucket is required when metadata type is s3 (ZT_METADATA_S3_BUCKET)")
		}
		if cfg.Metadata.S3Region == "" {
			cfg.Metadata.S3Region = "us-east-1"
		}
	default:
		return fmt.Errorf("invalid metadata type %q (must be filesystem or s3)", cfg.Metadata.Type)
	}

	return nil
}

// MetadataFilesystemPath returns the resolved filesystem root for metadata storage.
func (cfg *Config) MetadataFilesystemPath() string {
	if cfg.Metadata.Type != "filesystem" {
		return ""
	}
	return cfg.Metadata.Path
}

// MetadataPathDisplay returns the path shown in admin (filesystem path or S3 prefix).
func (cfg *Config) MetadataPathDisplay() string {
	if cfg.Metadata.Type == "filesystem" {
		return cfg.Metadata.Path
	}
	return cfg.Metadata.S3Prefix
}
