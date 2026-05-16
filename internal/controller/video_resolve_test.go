package controller

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
)

func TestResolveVideoFile_CrossLibraryTriage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defaultRoot := t.TempDir()
	otherRoot := t.TempDir()
	_ = os.MkdirAll(filepath.Join(defaultRoot, "triage"), 0o755)
	videoFile := filepath.Join(defaultRoot, "triage", "test.mp4")
	if err := os.WriteFile(videoFile, []byte("video"), 0o644); err != nil {
		t.Fatal(err)
	}

	ctrl := setupVideoController(t)
	var defaultLib model.Library
	if err := ctrl.datastore.First(&defaultLib).Error; err != nil {
		t.Fatal(err)
	}
	defaultLib.Config = model.LibraryConfig{Filesystem: &model.LibraryConfigFilesystem{Path: defaultRoot}}
	if err := ctrl.datastore.Save(&defaultLib).Error; err != nil {
		t.Fatal(err)
	}
	ctrl.storageResolver.Invalidate(defaultLib.ID)
	otherLib := model.Library{
		Name: "Other",
		Type: model.LibraryTypeFilesystem,
		Config: model.LibraryConfig{
			Filesystem: &model.LibraryConfigFilesystem{Path: otherRoot},
		},
	}
	if err := ctrl.datastore.Create(&otherLib).Error; err != nil {
		t.Fatal(err)
	}

	wrongPath := "/videos/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa/video.mp4"
	vid := &model.Video{
		ID:        "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		Name:      "Broken",
		Filename:  "test.mp4",
		Type:      "v",
		Imported:  true,
		Path:      &wrongPath,
		LibraryID: &otherLib.ID,
	}
	if err := ctrl.datastore.Create(vid).Error; err != nil {
		t.Fatal(err)
	}

	ctrl.config = &config.Config{DefaultLibraryID: defaultLib.ID}
	ctrl.storageResolver = storage.NewResolver(ctrl.datastore)

	store, path, found := ctrl.resolveVideoFile(vid)
	if !found {
		t.Fatal("expected to find video on default library triage path")
	}
	if path != "triage/test.mp4" {
		t.Fatalf("path = %q, want triage/test.mp4", path)
	}
	fs, ok := store.(*storage.Filesystem)
	if !ok {
		t.Fatalf("expected filesystem store, got %T", store)
	}
	if fs.FullPath(path) != videoFile {
		t.Fatalf("full path = %q, want %q", fs.FullPath(path), videoFile)
	}
}
