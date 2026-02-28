package model

import (
	"gorm.io/gorm"
)

// EnsureDefaultLibrary creates a default filesystem library from mediaPath if no libraries exist.
// Returns the default library ID (existing or newly created) and any error.
func EnsureDefaultLibrary(db *gorm.DB, mediaPath string) (string, error) {
	var count int64
	err := db.Model(&Library{}).Count(&count).Error
	if err != nil {
		return "", err
	}
	if count > 0 {
		var def Library
		err = db.Where("is_default = ?", true).First(&def).Error
		if err != nil {
			err = db.First(&def).Error
			if err != nil {
				return "", err
			}
		}
		return def.ID, nil
	}
	// Use fixed UUID so migration default (videos.library_id) points to this library without backfill.
	const defaultLibraryUUID = "00000000-0000-0000-0000-000000000000"
	lib := Library{
		ID:        defaultLibraryUUID,
		Name:      "Default",
		Type:      LibraryTypeFilesystem,
		IsDefault: true,
		Config: LibraryConfig{
			Filesystem: &LibraryConfigFilesystem{Path: mediaPath},
		},
	}
	if err := db.Create(&lib).Error; err != nil {
		return "", err
	}
	return lib.ID, nil
}

// BackfillVideoLibraryID sets library_id to defaultLibraryID for all videos where library_id is null.
func BackfillVideoLibraryID(db *gorm.DB, defaultLibraryID string) error {
	return db.Model(&Video{}).Where("library_id IS NULL").Updates(map[string]any{"library_id": defaultLibraryID}).Error
}
