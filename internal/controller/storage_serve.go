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
	}
	g.DataFromReader(http.StatusOK, -1, contentType, rc, nil)
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
	store, err := c.storageResolver.Storage(c.videoLibraryID(video))
	if err != nil {
		return ""
	}
	ps, ok := store.(storage.PreviewableStorage)
	if !ok {
		return ""
	}
	var path string
	if video.Imported {
		path = video.RelativePath()
	} else {
		path = filepath.Join("triage", video.Filename)
	}
	url, err := ps.PresignGet(g.Request.Context(), path, time.Hour)
	if err != nil || url == "" {
		return ""
	}
	return url
}
