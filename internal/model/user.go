package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Username  string
	Password  string
	Admin     bool
	Session   UserSession
}

// UUID pre-hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
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
