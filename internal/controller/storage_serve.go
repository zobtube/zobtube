package controller

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
)

// serveFromStorage serves the object at path from the given storage (filesystem or S3).
func (c *Controller) serveFromStorage(g *gin.Context, store storage.Storage, path string) {
	if store == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "storage not available"})
		return
	}
	if fs, ok := store.(*storage.Filesystem); ok {
		g.File(fs.FullPath(path))
		return
	}
	rc, err := store.Open(path)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rc.Close()
	contentType := "application/octet-stream"
	switch filepath.Ext(path) {
	case ".mp4", ".webm", ".mkv":
		contentType = "video/mp4"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	}
	g.DataFromReader(http.StatusOK, -1, contentType, rc, nil)
}

// resolveVideoFile locates the video blob on the assigned library, then on any
// other library (recovery after a partial library migration).
func (c *Controller) resolveVideoFile(video *model.Video) (storage.Storage, string, bool) {
	candidates := video.StoragePathCandidates()
	if !video.Imported && video.Filename != "" {
		candidates = []string{filepath.Join("triage", video.Filename)}
	}
	libID := c.videoLibraryID(video)
	store, err := c.storageResolver.Storage(libID)
	if err == nil {
		if p, ok, _ := storage.FirstExistingPath(store, candidates); ok {
			return store, p, true
		}
	}
	var libs []model.Library
	if c.datastore.Find(&libs).Error != nil {
		if store != nil {
			return store, video.RelativePath(), false
		}
		return nil, "", false
	}
	for _, lib := range libs {
		if lib.ID == libID {
			continue
		}
		alt, err := c.storageResolver.Storage(lib.ID)
		if err != nil {
			continue
		}
		if p, ok, _ := storage.FirstExistingPath(alt, candidates); ok {
			return alt, p, true
		}
	}
	if store != nil {
		if !video.Imported && video.Filename != "" {
			return store, filepath.Join("triage", video.Filename), false
		}
		return store, video.RelativePath(), false
	}
	return nil, "", false
}

// videoLibraryID returns the library ID for the video (or default if unset).
func (c *Controller) videoLibraryID(video *model.Video) string {
	if video.LibraryID != nil && *video.LibraryID != "" {
		return *video.LibraryID
	}
	return c.config.DefaultLibraryID
}

// videoStreamURL returns a direct stream URL when the video's storage supports it (e.g. S3 presigned).
// Otherwise returns empty string; frontend falls back to /api/video/:id/stream.
func (c *Controller) videoStreamURL(g *gin.Context, video *model.Video) string {
	if c.storageResolver == nil {
		return ""
	}
	store, path, _ := c.resolveVideoFile(video)
	if store == nil {
		return ""
	}
	ps, ok := store.(storage.PreviewableStorage)
	if !ok {
		return ""
	}
	url, err := ps.PresignGet(g.Request.Context(), path, time.Hour)
	if err != nil || url == "" {
		return ""
	}
	return url
}
