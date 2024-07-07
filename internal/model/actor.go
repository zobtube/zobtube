package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Actor struct {
	ID        string `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	Videos    []Video `gorm:"many2many:video_actors;"`
	Thumbnail bool
	Sex       string `gorm:"size:1;"`
	Aliases   []ActorAlias
	Links     []ActorLink
}

// UUID pre-hook
func (a *Actor) BeforeCreate(tx *gorm.DB) error {
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

var sexTypesAsString = map[string]string{
	"m": "male",
	"f": "female",
	"s": "shemale",
}

func (a *Actor) SexTypeAsString() string {
	return sexTypesAsString[a.Sex]
}

func (a *Actor) AliasesAsNiceString() string {
	var aliases []string

	for _, alias := range a.Aliases {
		aliases = append(aliases, alias.Name)
	}
	return strings.Join(aliases, " / ")
}

// URLs
func (a *Actor) URLView() string {
	return fmt.Sprintf("/actor/%s", a.ID)
}

func (a *Actor) URLThumb() string {
	return fmt.Sprintf("/actor/%s/thumb", a.ID)
}
