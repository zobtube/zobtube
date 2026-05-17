package controller

import (
	"mime"
	"path/filepath"
	"strings"

	"github.com/zobtube/zobtube/internal/model"
)

var allowedPhotoExtensions = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
}

func photoMimeForFilename(name string) (string, bool) {
	ext := strings.ToLower(filepath.Ext(name))
	m, ok := allowedPhotoExtensions[ext]
	return m, ok
}

func sanitizePhotoBasename(name string) string {
	base := filepath.Base(name)
	base = strings.ReplaceAll(base, "\\", "_")
	base = strings.ReplaceAll(base, "/", "_")
	if base == "" || base == "." || base == ".." {
		return ""
	}
	return base
}

func (c *Controller) photosetLibraryID(ps *model.Photoset) string {
	if ps.LibraryID != nil && *ps.LibraryID != "" {
		return *ps.LibraryID
	}
	return c.config.DefaultLibraryID
}

func mergeActors(base []model.Actor, extra []model.Actor) []model.Actor {
	seen := make(map[string]struct{})
	out := make([]model.Actor, 0, len(base)+len(extra))
	for _, a := range base {
		if _, ok := seen[a.ID]; ok {
			continue
		}
		seen[a.ID] = struct{}{}
		out = append(out, a)
	}
	for _, a := range extra {
		if _, ok := seen[a.ID]; ok {
			continue
		}
		seen[a.ID] = struct{}{}
		out = append(out, a)
	}
	return out
}

func mergeCategoryMap(base []model.CategorySub, extra []model.CategorySub) map[string]string {
	out := make(map[string]string)
	for _, cat := range base {
		out[cat.ID] = cat.Name
	}
	for _, cat := range extra {
		out[cat.ID] = cat.Name
	}
	return out
}

func effectiveChannel(ps *model.Photoset, photo *model.Photo) *model.Channel {
	if photo.Channel != nil {
		return photo.Channel
	}
	return ps.Channel
}

func detectContentType(path, fallback string) string {
	if fallback != "" {
		return fallback
	}
	if ext := strings.ToLower(filepath.Ext(path)); ext != "" {
		if ct := mime.TypeByExtension(ext); ct != "" {
			return ct
		}
	}
	return "application/octet-stream"
}
