package model

import (
	"path/filepath"

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

// EnsureDefaultOrganization creates the bootstrap "v1" Organization if no
// Organization exists yet. The default template reproduces the legacy
// hardcoded layout (<typePlural>/<id>/video.mp4) so existing imported videos
// can be associated to it without moving any files. Returns the active
// organization ID (existing active one, or the newly created default).
func EnsureDefaultOrganization(db *gorm.DB) (string, error) {
	var count int64
	if err := db.Model(&Organization{}).Count(&count).Error; err != nil {
		return "", err
	}
	if count > 0 {
		var active Organization
		err := db.Where("active = ?", true).First(&active).Error
		if err == nil {
			return active.ID, nil
		}
		var any Organization
		if err := db.First(&any).Error; err != nil {
			return "", err
		}
		return any.ID, nil
	}
	org := Organization{
		ID:       DefaultOrganizationUUID,
		Name:     "Default (legacy layout)",
		Template: DefaultOrganizationTemplate,
		Active:   true,
	}
	if err := db.Create(&org).Error; err != nil {
		return "", err
	}
	return org.ID, nil
}

// BackfillVideoOrganization assigns defaultOrganizationID to every imported
// video that has no organization yet, and sets the resolved Path column from
// the legacy hardcoded layout. This is safe to run on every boot because it
// only touches rows where organization_id IS NULL.
func BackfillVideoOrganization(db *gorm.DB, defaultOrganizationID string) error {
	var videos []Video
	err := db.Where("imported = ? AND organization_id IS NULL", true).Find(&videos).Error
	if err != nil {
		return err
	}
	for i := range videos {
		v := &videos[i]
		legacyPath := filepath.Join(v.FolderRelativePath(), "video.mp4")
		orgID := defaultOrganizationID
		v.OrganizationID = &orgID
		v.Path = &legacyPath
		if err := db.Model(v).Updates(map[string]any{
			"organization_id": orgID,
			"path":            legacyPath,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
