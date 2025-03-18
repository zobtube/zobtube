package model

type VideoView struct {
	VideoID string `gorm:"primaryKey"`
	Video   Video
	UserID  string `gorm:"primaryKey"`
	User    User
	Count   int
}
