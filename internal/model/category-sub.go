package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategorySub struct {
	ID        string `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	Category  string `gorm:"type:uuid"`
	Thumbnail bool
}

// UUID pre-hook
func (c *CategorySub) BeforeCreate(tx *gorm.DB) error {
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

func (c *CategorySub) URLThumb() string {
	return fmt.Sprintf("/category-sub/%s/thumb", c.ID)
}
