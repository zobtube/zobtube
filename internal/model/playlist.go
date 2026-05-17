package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Playlist struct {
	ID        string    `gorm:"type:uuid;primary_key"`
	UserID    string    `gorm:"type:uuid;index;not null"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Playlist) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return nil
}

type PlaylistVideo struct {
	PlaylistID string    `gorm:"type:uuid;primaryKey"`
	VideoID    string    `gorm:"type:uuid;primaryKey"`
	Position   int       `gorm:"not null;default:0"`
	AddedAt    time.Time `gorm:"not null"`
}
