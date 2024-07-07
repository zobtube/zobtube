package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActorLink struct {
	ID        string `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Provider  string
	ActorID   string
	Actor     Actor
	URL       string
}

// UUID pre-hook
func (a *ActorLink) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "00000000-0000-0000-0000-000000000000" {
		a.ID = uuid.NewString()
		return nil
	}

	if a.ID == "" {
		a.ID = uuid.NewString()
		return nil
	}

	return nil
}
