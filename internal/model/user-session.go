package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSession struct {
	ID         string `gorm:"type:uuid;primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserID     string
	ValidUntil time.Time
}

// UUID pre-hook
func (u *UserSession) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "00000000-0000-0000-0000-000000000000" {
		u.ID = uuid.NewString()
		return nil
	}

	if u.ID == "" {
		u.ID = uuid.NewString()
		return nil
	}

	return nil
}
