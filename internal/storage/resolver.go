package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/gorm"

	"github.com/zobtube/zobtube/internal/model"
)

// Resolver returns Storage for a library by ID (cached).
type Resolver struct {
	db    *gorm.DB
	cache sync.Map // libraryID -> Storage
}

// NewResolver returns a new storage resolver.
func NewResolver(db *gorm.DB) *Resolver {
	return &Resolver{db: db}
}

// Storage returns the Storage for the given library ID. Results are cached.
func (r *Resolver) Storage(libraryID string) (Storage, error) {
	if v, ok := r.cache.Load(libraryID); ok {
		return v.(Storage), nil
	}
	store, err := r.load(libraryID)
	if err != nil {
		return nil, err
	}
	r.cache.Store(libraryID, store)
	return store, nil
}

// Invalidate removes the cached Storage for the given library (e.g. after config change).
func (r *Resolver) Invalidate(libraryID string) {
	r.cache.Delete(libraryID)
}

func (r *Resolver) load(libraryID string) (Storage, error) {
	var lib model.Library
	err := r.db.First(&lib, "id = ?", libraryID).Error
	if err != nil {
		return nil, err
	}
	switch lib.Type {
	case model.LibraryTypeFilesystem:
		path := ""
		if lib.Config.Filesystem != nil {
			path = lib.Config.Filesystem.Path
		}
		return NewFilesystem(path), nil
	case model.LibraryTypeS3:
		if lib.Config.S3 == nil {
			return nil, ErrS3ConfigMissing
		}
		cfg := lib.Config.S3
		client, err := newS3Client(cfg.Region, cfg.Endpoint, nil)
		if err != nil {
			return nil, err
		}
		return NewS3(client, cfg.Bucket, cfg.Prefix), nil
	default:
		return nil, ErrUnknownLibraryType
	}
}

var (
	ErrS3ConfigMissing    = errors.New("library: s3 config missing")
	ErrUnknownLibraryType = errors.New("library: unknown type")
)

// newS3Client builds an S3 client. Credentials come from env (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY) or IAM.
// endpoint is optional (e.g. for Minio); creds are optional (for static Minio credentials).
func newS3Client(region, endpoint string, creds *struct{ AccessKey, SecretKey string }) (*s3.Client, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}
	if endpoint != "" {
		opts = append(opts, config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               endpoint,
					SigningRegion:     region,
					HostnameImmutable: true,
				}, nil
			},
		)))
	}
	if creds != nil && creds.AccessKey != "" && creds.SecretKey != "" {
		opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			creds.AccessKey, creds.SecretKey, "",
		)))
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(cfg), nil
}
