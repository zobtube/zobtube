package model

type Configuration struct {
	ID                 int `gorm:"type:primary_key"`
	UserAuthentication bool
}
