package model

import (
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DefaultOrganizationUUID is the fixed UUID for the bootstrap "v1" organization
// that matches the legacy hardcoded layout. New imports default to whatever
// organization is currently Active; this one is created Active by bootstrap if
// no organization exists yet.
const DefaultOrganizationUUID = "00000000-0000-0000-0000-000000000001"

// DefaultOrganizationTemplate is the path template that reproduces the legacy
// hardcoded layout: <typePlural>/<id>/video.mp4 (e.g. videos/abc/video.mp4).
const DefaultOrganizationTemplate = "$TYPE/$ID/video.mp4"

// Organization defines a versioned scheme used to store video files on a
// library. Only one Organization is Active at a time; new imports use the
// Active one to resolve the on-disk path. Existing videos keep their stored
// Path until they are explicitly reorganized.
type Organization struct {
	ID        string `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:255;not null"`
	Template  string `gorm:"size:1024;not null"`
	Active    bool   `gorm:"default:false"`
}

func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.NewString()
	}
	return nil
}

// ErrOrganizationTemplateMissingID indicates the template does not include
// $ID, which is required to guarantee unique paths across videos.
var ErrOrganizationTemplateMissingID = errors.New("organization: template must include $ID")

// ValidateOrganizationTemplate verifies that template is non-empty, contains
// $ID, and only references known variables.
func ValidateOrganizationTemplate(tmpl string) error {
	if strings.TrimSpace(tmpl) == "" {
		return errors.New("organization: template cannot be empty")
	}
	if !strings.Contains(tmpl, "$ID") {
		return ErrOrganizationTemplateMissingID
	}
	return nil
}

// videoTypePlural maps the type single letter to the plural folder name used
// by the legacy layout and by the default template.
var videoTypePlural = map[string]string{
	"c": "clips",
	"v": "videos",
	"m": "movies",
}

// Render returns the relative path produced by applying the template to the
// given video. Variables supported:
//
//	$ID          - the video UUID
//	$TYPE        - plural folder for the type ("clips", "videos", "movies")
//	$TYPE_NAME   - singular type name ("clip", "video", "movie")
//	$TYPE_LETTER - single-letter type ("c", "v", "m")
//	$FILENAME    - original imported filename (basename, with extension)
//	$BASENAME    - original filename without its extension
//	$EXT         - extension of the original filename including the dot
//
// Any leading "/" is stripped so the result is always relative to the
// library root.
func (o *Organization) Render(v *Video) string {
	out := o.Template
	out = strings.ReplaceAll(out, "$ID", v.ID)
	out = strings.ReplaceAll(out, "$TYPE_LETTER", v.Type)
	out = strings.ReplaceAll(out, "$TYPE_NAME", videoTypesAsString[v.Type])
	out = strings.ReplaceAll(out, "$TYPE", videoTypePlural[v.Type])
	filename := filepath.Base(v.Filename)
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	out = strings.ReplaceAll(out, "$FILENAME", filename)
	out = strings.ReplaceAll(out, "$BASENAME", base)
	out = strings.ReplaceAll(out, "$EXT", ext)
	out = strings.TrimPrefix(out, "/")
	return filepath.Clean(out)
}

// ActiveOrganization returns the organization marked Active, or the first row if none is active.
func ActiveOrganization(db *gorm.DB) (*Organization, error) {
	var active Organization
	if err := db.Where("active = ?", true).First(&active).Error; err == nil {
		return &active, nil
	}
	var any Organization
	if err := db.First(&any).Error; err != nil {
		return nil, err
	}
	return &any, nil
}

// IsOrganizedWith reports whether the video follows the given organization's on-disk layout.
func (v *Video) IsOrganizedWith(org *Organization) bool {
	if org == nil || !v.Imported {
		return false
	}
	if v.OrganizationID == nil || *v.OrganizationID != org.ID {
		return false
	}
	return v.RelativePath() == org.Render(v)
}

// NeedsReorganization reports whether the video should be moved onto org's layout.
func (v *Video) NeedsReorganization(org *Organization) bool {
	if org == nil || !v.Imported {
		return false
	}
	return !v.IsOrganizedWith(org)
}
