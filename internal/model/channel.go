package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Channel struct {
	ID        string `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	Thumbnail bool
}

// UUID pre-hook
func (c *Channel) BeforeCreate(tx *gorm.DB) error {
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

func (c *Channel) URLView() string {
	return fmt.Sprintf("/channel/%s", c.ID)
}

func (c *Channel) URLThumb() string {
	return fmt.Sprintf("/api/channel/%s/thumb", c.ID)
}

func (c *Channel) URLAdmEdit() string {
	return fmt.Sprintf("/channel/%s/edit", c.ID)
}
