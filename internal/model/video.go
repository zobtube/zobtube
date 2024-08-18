package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Video model defines the generic video type used for videos, clips and movies
type Video struct {
	ID            string `gorm:"type:uuid;primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Name          string
	Filename      string
	Actors        []Actor  `gorm:"many2many:video_actors;"`
	Channel       *Channel `gorm:"foreignKey:ID"`
	Thumbnail     bool
	ThumbnailMini bool
	Duration      time.Duration
	Type          string `gorm:"size:1;"`
	Imported      bool   `gorm:"default:false"`
}

var videoTypesAsString = map[string]string{
	"c": "clip",
	"v": "video",
	"m": "movie",
}

func (v *Video) BeforeCreate(tx *gorm.DB) error {
	if v.ID == "00000000-0000-0000-0000-000000000000" {
		v.ID = uuid.NewString()
		return nil
	}

	if v.ID == "" {
		v.ID = uuid.NewString()
		return nil
	}
	return nil
}

func (v *Video) TypeAsString() string {
	return videoTypesAsString[v.Type]
}

func (v *Video) URLView() string {
	return fmt.Sprintf("/video/%s", v.ID)
}

func (v *Video) URLThumb() string {
	return fmt.Sprintf("/video/%s/thumb", v.ID)
}

func (v *Video) URLThumbXS() string {
	return fmt.Sprintf("/video/%s/thumb_xs", v.ID)
}

func (v *Video) URLStream() string {
	return fmt.Sprintf("/video/%s/stream", v.ID)
}

func (v *Video) URLAdmEdit() string {
	return fmt.Sprintf("/video/%s/edit", v.ID)
}

func (v *Video) HasDuration() bool {
	return true
}

func (v *Video) NiceDuration() string {
	d := v.Duration
	d = d.Round(time.Second)

	h := d / time.Hour
	d -= h * time.Hour

	m := d / time.Minute
	d -= m * time.Minute

	s := d / time.Second
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	} else {
		return fmt.Sprintf("%02d:%02d", m, s)
	}
}

func (v *Video) String() string {
	return v.Name
}
