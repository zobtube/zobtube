package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LibraryType is the storage backend type for a library.
type LibraryType string

const (
	LibraryTypeFilesystem LibraryType = "filesystem"
	LibraryTypeS3         LibraryType = "s3"
)

// LibraryConfigFilesystem is the config for a filesystem library.
type LibraryConfigFilesystem struct {
	Path string `json:"path"`
}

// LibraryConfigS3 is the config for an S3 library.
// AccessKeyID and SecretAccessKey are optional; if not set, the default credential chain is used (env AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, or IAM).
type LibraryConfigS3 struct {
	Bucket          string `json:"bucket"`
	Region          string `json:"region,omitempty"`
	Prefix          string `json:"prefix,omitempty"`
	Endpoint        string `json:"endpoint,omitempty"` // for Minio etc.
	AccessKeyID     string `json:"access_key_id,omitempty"`
	SecretAccessKey string `json:"secret_access_key,omitempty"`
}

// LibraryConfig holds either filesystem or S3 config as JSON.
type LibraryConfig struct {
	Filesystem *LibraryConfigFilesystem `json:"filesystem,omitempty"`
	S3         *LibraryConfigS3         `json:"s3,omitempty"`
}

// Value implements driver.Valuer for GORM.
func (c LibraryConfig) Value() (driver.Value, error) {
	if c.Filesystem == nil && c.S3 == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

// Scan implements sql.Scanner for GORM.
func (c *LibraryConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("library config: invalid type")
	}
	return json.Unmarshal(b, c)
}

// Library is a named media storage target (filesystem or S3).
type Library struct {
	ID        string `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string        `gorm:"size:255;not null"`
	Type      LibraryType   `gorm:"size:32;not null"`
	Config    LibraryConfig `gorm:"type:text"`     // JSON
	IsDefault bool          `gorm:"default:false"` // one library is used for actor/channel/category assets
}

func (l *Library) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.NewString()
	}
	return nil
}
