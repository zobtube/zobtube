package video

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
	"github.com/zobtube/zobtube/internal/task/common"
)

func setupTaskContext(t *testing.T, libRoot string) *common.Context {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(
		&model.Library{}, &model.Organization{}, &model.Configuration{}, &model.Video{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	libID, err := model.EnsureDefaultLibrary(db, libRoot)
	if err != nil {
		t.Fatalf("ensure library: %v", err)
	}
	orgID, err := model.EnsureDefaultOrganization(db)
	if err != nil {
		t.Fatalf("ensure organization: %v", err)
	}
	_ = orgID
	return &common.Context{
		DB:              db,
		Config:          &config.Config{DefaultLibraryID: libID},
		StorageResolver: storage.NewResolver(db),
		MetadataStorage: storage.NewFilesystem(libRoot),
	}
}

func TestImportFromTriage_MovesFileToActiveOrgPath(t *testing.T) {
	tmp := t.TempDir()
	ctx := setupTaskContext(t, tmp)

	v := &model.Video{ID: "11111111-1111-1111-1111-111111111111", Name: "raw", Filename: "raw.mp4", Type: "v"}
	if err := ctx.DB.Create(v).Error; err != nil {
		t.Fatal(err)
	}
	libStore := storage.NewFilesystem(tmp)
	_ = libStore.MkdirAll("triage")
	wc, err := libStore.Create(filepath.Join("triage", "raw.mp4"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := wc.Write([]byte("payload")); err != nil {
		t.Fatal(err)
	}
	_ = wc.Close()

	if msg, err := importFromTriage(ctx, common.Parameters{"videoID": v.ID}); err != nil {
		t.Fatalf("importFromTriage failed: %v (%s)", err, msg)
	}

	var got model.Video
	if err := ctx.DB.First(&got, "id = ?", v.ID).Error; err != nil {
		t.Fatal(err)
	}
	if !got.Imported {
		t.Error("expected Imported=true")
	}
	if got.OrganizationID == nil {
		t.Fatal("expected organization_id to be set")
	}
	wantPath := filepath.Join("videos", v.ID, "video.mp4")
	if got.Path == nil || *got.Path != wantPath {
		t.Fatalf("expected stored path %q, got %v", wantPath, got.Path)
	}
	if _, err := os.Stat(filepath.Join(tmp, wantPath)); err != nil {
		t.Errorf("destination file missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "triage", "raw.mp4")); !os.IsNotExist(err) {
		t.Errorf("triage file should be removed, stat err = %v", err)
	}
}

func TestImportFromTriage_SkipReorganization_KeepsFileInTriage(t *testing.T) {
	tmp := t.TempDir()
	ctx := setupTaskContext(t, tmp)

	v := &model.Video{ID: "22222222-2222-2222-2222-222222222222", Name: "raw", Filename: "raw.mp4", Type: "v"}
	if err := ctx.DB.Create(v).Error; err != nil {
		t.Fatal(err)
	}
	libStore := storage.NewFilesystem(tmp)
	_ = libStore.MkdirAll("triage")
	wc, _ := libStore.Create(filepath.Join("triage", "raw.mp4"))
	_, _ = wc.Write([]byte("payload"))
	_ = wc.Close()

	if msg, err := importFromTriage(ctx, common.Parameters{"videoID": v.ID, "skipReorganization": "true"}); err != nil {
		t.Fatalf("importFromTriage failed: %v (%s)", err, msg)
	}

	var got model.Video
	if err := ctx.DB.First(&got, "id = ?", v.ID).Error; err != nil {
		t.Fatal(err)
	}
	if !got.Imported {
		t.Error("expected Imported=true even when skipping reorg")
	}
	if got.OrganizationID != nil {
		t.Errorf("expected nil organization_id, got %v", got.OrganizationID)
	}
	wantPath := filepath.Join("triage", "raw.mp4")
	if got.Path == nil || *got.Path != wantPath {
		t.Fatalf("expected stored path %q, got %v", wantPath, got.Path)
	}
	if _, err := os.Stat(filepath.Join(tmp, "triage", "raw.mp4")); err != nil {
		t.Errorf("triage file should still exist: %v", err)
	}
}

func TestImportFromTriage_HonorsGlobalConfigDefault(t *testing.T) {
	tmp := t.TempDir()
	ctx := setupTaskContext(t, tmp)
	if err := ctx.DB.Create(&model.Configuration{ID: 1}).Error; err != nil {
		t.Fatal(err)
	}
	// Force ReorganizeOnImport=false; GORM default kicks in on bool zero-values
	// during INSERT so a plain struct assignment isn't enough.
	if err := ctx.DB.Model(&model.Configuration{}).Where("id = ?", 1).Update("reorganize_on_import", false).Error; err != nil {
		t.Fatal(err)
	}

	v := &model.Video{ID: "33333333-3333-3333-3333-333333333333", Name: "raw", Filename: "raw.mp4", Type: "v"}
	if err := ctx.DB.Create(v).Error; err != nil {
		t.Fatal(err)
	}
	libStore := storage.NewFilesystem(tmp)
	_ = libStore.MkdirAll("triage")
	wc, _ := libStore.Create(filepath.Join("triage", "raw.mp4"))
	_, _ = wc.Write([]byte("payload"))
	_ = wc.Close()

	if msg, err := importFromTriage(ctx, common.Parameters{"videoID": v.ID}); err != nil {
		t.Fatalf("importFromTriage failed: %v (%s)", err, msg)
	}

	var got model.Video
	_ = ctx.DB.First(&got, "id = ?", v.ID).Error
	if got.OrganizationID != nil {
		t.Errorf("expected reorganization skipped via global default, got organization_id %v", got.OrganizationID)
	}
	if got.Path == nil || *got.Path != filepath.Join("triage", "raw.mp4") {
		t.Errorf("expected path to triage, got %v", got.Path)
	}
}
