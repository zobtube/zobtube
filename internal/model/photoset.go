package model

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PhotosetStatus string

const (
	PhotosetStatusCreating PhotosetStatus = "creating"
	PhotosetStatusReady    PhotosetStatus = "ready"
	PhotosetStatusDeleting PhotosetStatus = "deleting"
)

// Photoset is an album composed of several photos.
type Photoset struct {
	ID             string `gorm:"type:uuid;primary_key"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	Name           string
	LibraryID      *string       `gorm:"type:uuid;index;default:00000000-0000-0000-0000-000000000000"`
	OrganizationID *string       `gorm:"type:uuid;index"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID;references:ID"`
	ChannelID      *string
	Channel        *Channel
	Actors         []Actor       `gorm:"many2many:photoset_actors;"`
	Categories     []CategorySub `gorm:"many2many:photoset_categories;"`
	Photos         []Photo
	CoverPhotoID   *string        `gorm:"type:uuid;index"`
	Status         PhotosetStatus `gorm:"default:creating"`
	Imported       bool           `gorm:"default:false"`
}

func (ps *Photoset) BeforeCreate(tx *gorm.DB) error {
	if ps.ID == "" {
		ps.ID = uuid.NewString()
	}
	return nil
}

func (ps *Photoset) FolderRelativePath() string {
	return filepath.Join("photosets", ps.ID)
}

func (ps *Photoset) URLView() string {
	return fmt.Sprintf("/photoset/%s", ps.ID)
}

func (ps *Photoset) URLCover() string {
	return fmt.Sprintf("/api/photoset/%s/cover", ps.ID)
}

func (ps *Photoset) URLAdmEdit() string {
	return fmt.Sprintf("/photoset/%s/edit", ps.ID)
}

func (ps *Photoset) String() string {
	return ps.Name
}

// IsOrganizedWith reports whether every imported photo follows the organization's layout.
func (ps *Photoset) IsOrganizedWith(org *Organization) bool {
	if org == nil || !ps.Imported {
		return false
	}
	if ps.OrganizationID == nil || *ps.OrganizationID != org.ID {
		return false
	}
	for i := range ps.Photos {
		p := &ps.Photos[i]
		if p.RelativePath(ps) != org.RenderPhotoset(ps, p) {
			return false
		}
	}
	return true
}

// NeedsReorganization reports whether the photoset should be moved onto org's layout.
func (ps *Photoset) NeedsReorganization(org *Organization) bool {
	if org == nil || !ps.Imported {
		return false
	}
	return !ps.IsOrganizedWith(org)
}
