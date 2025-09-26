package model

import (
	"time"

	"gorm.io/gorm"
)

type Provider struct {
	ID                  string `gorm:"type:string;primary_key"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
	NiceName            string
	Enabled             bool `gorm:"default:true"`
	AbleToSearchActor   bool `gorm:"default:false"`
	AbleToScrapePicture bool `gorm:"default:false"`
}
