package storage

import (
	"fmt"

	"github.com/zobtube/zobtube/internal/config"
)

// OpenMetadata builds a Storage backend from global metadata configuration.
func OpenMetadata(cfg *config.Config) (Storage, error) {
	if cfg == nil {
		return nil, fmt.Errorf("metadata storage: config is nil")
	}
	switch cfg.Metadata.Type {
	case "filesystem":
		return NewFilesystem(cfg.Metadata.Path), nil
	case "s3":
		var creds *StaticS3Credentials
		if cfg.Metadata.S3AccessKeyID != "" && cfg.Metadata.S3SecretAccessKey != "" {
			creds = &StaticS3Credentials{
				AccessKey: cfg.Metadata.S3AccessKeyID,
				SecretKey: cfg.Metadata.S3SecretAccessKey,
			}
		}
		client, err := NewS3Client(cfg.Metadata.S3Region, cfg.Metadata.S3Endpoint, creds)
		if err != nil {
			return nil, err
		}
		return NewS3(client, cfg.Metadata.S3Bucket, cfg.Metadata.S3Prefix), nil
	default:
		return nil, fmt.Errorf("metadata storage: unknown type %q", cfg.Metadata.Type)
	}
}
