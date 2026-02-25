package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApiToken struct {
	ID        string    `gorm:"type:uuid;primary_key"`
	UserID    string    `gorm:"type:uuid;index;not null"`
	Name      string    `gorm:"not null"`
	TokenHash string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (a *ApiToken) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	return nil
}
