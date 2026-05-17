package photoset

import (
	"path/filepath"
	"strings"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

var allowedPhotoExtensions = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
}

func photosetLibraryID(ctx *common.Context, ps *model.Photoset) string {
	if ps.LibraryID != nil && *ps.LibraryID != "" {
		return *ps.LibraryID
	}
	return ctx.Config.DefaultLibraryID
}

func photoMimeForFilename(name string) (string, bool) {
	ext := strings.ToLower(filepath.Ext(name))
	m, ok := allowedPhotoExtensions[ext]
	return m, ok
}

func sanitizeBasename(name string) string {
	base := filepath.Base(name)
	base = strings.ReplaceAll(base, "\\", "_")
	base = strings.ReplaceAll(base, "/", "_")
	if base == "" || base == "." || base == ".." {
		return ""
	}
	return base
}
