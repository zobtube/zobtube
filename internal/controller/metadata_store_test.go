package controller

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
)

func setupMetadataStoreController(t *testing.T, libPath, otherLibPath, metaPath string) *Controller {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("db: %v", err)
	}
	if err := db.AutoMigrate(&model.Library{}, &model.Video{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	libID, err := model.EnsureDefaultLibrary(db, libPath)
	if err != nil {
		t.Fatalf("library: %v", err)
	}
	otherLib := model.Library{
		ID:   "11111111-1111-1111-1111-111111111111",
		Name: "Other",
		Type: model.LibraryTypeFilesystem,
		Config: model.LibraryConfig{
			Filesystem: &model.LibraryConfigFilesystem{Path: otherLibPath},
		},
	}
	if err := db.Create(&otherLib).Error; err != nil {
		t.Fatalf("create lib: %v", err)
	}

	log := zerolog.Nop()
	shutdown := make(chan int, 1)
	ctrl := New(shutdown).(*Controller)
	ctrl.LoggerRegister(&log)
	ctrl.DatabaseRegister(db)
	cfg := &config.Config{}
	cfg.DefaultLibraryID = libID
	ctrl.ConfigurationRegister(cfg)
	ctrl.StorageResolverRegister(storage.NewResolver(db))
	ctrl.MetadataStorageRegister(storage.NewFilesystem(metaPath))

	return ctrl
}

func TestController_metadataStore_unmigratedUsesDefaultLibrary(t *testing.T) {
	tmp := t.TempDir()
	libPath := filepath.Join(tmp, "lib")
	metaPath := filepath.Join(tmp, "meta")
	otherPath := filepath.Join(tmp, "other")
	for _, p := range []string{libPath, metaPath, otherPath} {
		_ = os.MkdirAll(p, 0o750)
	}
	ctrl := setupMetadataStoreController(t, libPath, otherPath, metaPath)

	store, err := ctrl.metadataStore(false)
	if err != nil {
		t.Fatal(err)
	}
	fs, ok := store.(*storage.Filesystem)
	if !ok || fs.Root != libPath {
		t.Fatalf("expected default library root %q, got %#v", libPath, store)
	}
}

func TestController_metadataStore_migratedUsesMetadata(t *testing.T) {
	tmp := t.TempDir()
	libPath := filepath.Join(tmp, "lib")
	metaPath := filepath.Join(tmp, "meta")
	otherPath := filepath.Join(tmp, "other")
	for _, p := range []string{libPath, metaPath, otherPath} {
		_ = os.MkdirAll(p, 0o750)
	}
	ctrl := setupMetadataStoreController(t, libPath, otherPath, metaPath)

	store, err := ctrl.metadataStore(true)
	if err != nil {
		t.Fatal(err)
	}
	fs, ok := store.(*storage.Filesystem)
	if !ok || fs.Root != metaPath {
		t.Fatalf("expected metadata root %q, got %#v", metaPath, store)
	}
}

func TestController_videoThumbnailStore_unmigratedUsesVideoLibrary(t *testing.T) {
	tmp := t.TempDir()
	libPath := filepath.Join(tmp, "lib")
	metaPath := filepath.Join(tmp, "meta")
	otherPath := filepath.Join(tmp, "other")
	for _, p := range []string{libPath, metaPath, otherPath} {
		_ = os.MkdirAll(p, 0o750)
	}
	ctrl := setupMetadataStoreController(t, libPath, otherPath, metaPath)

	otherID := "11111111-1111-1111-1111-111111111111"
	video := &model.Video{ID: "22222222-2222-2222-2222-222222222222", Name: "v", Filename: "v.mp4", Type: "v", LibraryID: &otherID}
	if err := ctrl.datastore.Create(video).Error; err != nil {
		t.Fatal(err)
	}

	store, err := ctrl.videoThumbnailStore(video)
	if err != nil {
		t.Fatal(err)
	}
	fs, ok := store.(*storage.Filesystem)
	if !ok || fs.Root != otherPath {
		t.Fatalf("expected video library root %q, got %#v", otherPath, store)
	}
}

func TestController_videoThumbnailStore_migratedUsesMetadata(t *testing.T) {
	tmp := t.TempDir()
	libPath := filepath.Join(tmp, "lib")
	metaPath := filepath.Join(tmp, "meta")
	otherPath := filepath.Join(tmp, "other")
	for _, p := range []string{libPath, metaPath, otherPath} {
		_ = os.MkdirAll(p, 0o750)
	}
	ctrl := setupMetadataStoreController(t, libPath, otherPath, metaPath)

	video := &model.Video{ID: "22222222-2222-2222-2222-222222222222", Name: "v", Filename: "v.mp4", Type: "v", Migrated: true}
	if err := ctrl.datastore.Create(video).Error; err != nil {
		t.Fatal(err)
	}

	store, err := ctrl.videoThumbnailStore(video)
	if err != nil {
		t.Fatal(err)
	}
	fs, ok := store.(*storage.Filesystem)
	if !ok || fs.Root != metaPath {
		t.Fatalf("expected metadata root %q, got %#v", metaPath, store)
	}
}
