package model

type Configuration struct {
	ID                  int `gorm:"type:primary_key"`
	UserAuthentication  bool
	OfflineMode         bool `gorm:"default:false"`
	ReorganizeOnImport  bool `gorm:"default:true"` // when true, new imports move files to the path defined by the active Organization
}
