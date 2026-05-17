package model

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Photo is a single image belonging to a photoset.
type Photo struct {
	ID            string `gorm:"type:uuid;primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	PhotosetID    string         `gorm:"type:uuid;index;not null"`
	Filename      string
	Path          *string `gorm:"size:1024"`
	Position      int
	Width         int
	Height        int
	SizeBytes     int64
	Mime          string `gorm:"size:32"`
	ThumbnailMini bool

	ChannelID  *string
	Channel    *Channel
	Actors     []Actor       `gorm:"many2many:photo_actors;"`
	Categories []CategorySub `gorm:"many2many:photo_categories;"`
}

func (p *Photo) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return nil
}

// RelativePath returns the relative path of the photo file on its library.
func (p *Photo) RelativePath(ps *Photoset) string {
	if p.Path != nil && *p.Path != "" {
		return *p.Path
	}
	return filepath.Join(ps.FolderRelativePath(), p.Filename)
}

// StoragePathCandidates returns relative paths to probe on a library.
func (p *Photo) StoragePathCandidates(ps *Photoset) []string {
	seen := make(map[string]struct{})
	var out []string
	add := func(path string) {
		if path == "" {
			return
		}
		path = filepath.Clean(path)
		if _, ok := seen[path]; ok {
			return
		}
		seen[path] = struct{}{}
		out = append(out, path)
	}
	if p.Path != nil && *p.Path != "" {
		add(*p.Path)
	}
	add(filepath.Join(ps.FolderRelativePath(), p.Filename))
	return out
}

func (p *Photo) ThumbnailMiniRelativePath(ps *Photoset) string {
	return filepath.Join("photosets", ps.ID, p.ID, "thumb-mini.jpg")
}

func (p *Photo) URLStream() string {
	return fmt.Sprintf("/api/photo/%s/stream", p.ID)
}

func (p *Photo) URLThumbMini() string {
	return fmt.Sprintf("/api/photo/%s/thumb_mini", p.ID)
}
