package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        string `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	Priority  int
	Sub       []CategorySub `gorm:"foreignKey:Category;references:ID"`
}

// UUID pre-hook
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "00000000-0000-0000-0000-000000000000" {
		c.ID = uuid.NewString()
		return nil
	}

	if c.ID == "" {
		c.ID = uuid.NewString()
		return nil
	}

	return nil
}
