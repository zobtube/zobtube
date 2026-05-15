package metamigrate

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

func testMigrateContext(t *testing.T, libRoot, metaRoot string) *common.Context {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.Library{}, &model.Actor{}, &model.Channel{}, &model.CategorySub{}, &model.Video{}); err != nil {
		t.Fatal(err)
	}
	libID, err := model.EnsureDefaultLibrary(db, libRoot)
	if err != nil {
		t.Fatal(err)
	}
	return &common.Context{
		DB:              db,
		Config:          &config.Config{DefaultLibraryID: libID},
		StorageResolver: storage.NewResolver(db),
		MetadataStorage: storage.NewFilesystem(metaRoot),
	}
}

func TestMigrateAll_actorThumbnail(t *testing.T) {
	tmp := t.TempDir()
	libRoot := filepath.Join(tmp, "lib")
	metaRoot := filepath.Join(tmp, "meta")
	_ = os.MkdirAll(libRoot, 0o750)
	_ = os.MkdirAll(metaRoot, 0o750)

	ctx := testMigrateContext(t, libRoot, metaRoot)
	libStore := storage.NewFilesystem(libRoot)

	actor := model.Actor{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "A", Thumbnail: true}
	if err := ctx.DB.Create(&actor).Error; err != nil {
		t.Fatal(err)
	}
	thumbPath := filepath.Join("actors", actor.ID, "thumb.jpg")
	if err := libStore.MkdirAll(filepath.Dir(thumbPath)); err != nil {
		t.Fatal(err)
	}
	f, err := libStore.Create(thumbPath)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.Write([]byte("jpeg"))
	_ = f.Close()

	msg, err := migrateAll(ctx, nil)
	if err != nil {
		t.Fatalf("migrateAll: %v (%s)", err, msg)
	}

	var updated model.Actor
	if err := ctx.DB.First(&updated, "id = ?", actor.ID).Error; err != nil {
		t.Fatal(err)
	}
	if !updated.Migrated {
		t.Fatal("expected actor migrated")
	}
	metaStore := storage.NewFilesystem(metaRoot)
	ok, _ := metaStore.Exists(thumbPath)
	if !ok {
		t.Fatal("thumbnail missing on metadata storage")
	}
	ok, _ = libStore.Exists(thumbPath)
	if ok {
		t.Fatal("legacy thumbnail should be removed")
	}
	if _, err := os.Stat(filepath.Join(libRoot, "actors")); !os.IsNotExist(err) {
		t.Fatal("expected legacy actors folder removed when empty")
	}
}

func TestMigrateAll_categoryThumbnail(t *testing.T) {
	tmp := t.TempDir()
	libRoot := filepath.Join(tmp, "lib")
	metaRoot := filepath.Join(tmp, "meta")
	_ = os.MkdirAll(libRoot, 0o750)
	_ = os.MkdirAll(metaRoot, 0o750)

	ctx := testMigrateContext(t, libRoot, metaRoot)
	libStore := storage.NewFilesystem(libRoot)

	sub := model.CategorySub{ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb", Name: "S", Thumbnail: true}
	if err := ctx.DB.Create(&sub).Error; err != nil {
		t.Fatal(err)
	}
	thumbPath := filepath.Join("categories", sub.ID+".jpg")
	if err := libStore.MkdirAll("categories"); err != nil {
		t.Fatal(err)
	}
	f, err := libStore.Create(thumbPath)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.Write([]byte("jpeg"))
	_ = f.Close()

	msg, err := migrateAll(ctx, nil)
	if err != nil {
		t.Fatalf("migrateAll: %v (%s)", err, msg)
	}
	if _, err := os.Stat(filepath.Join(libRoot, "categories")); !os.IsNotExist(err) {
		t.Fatal("expected legacy categories folder removed when empty")
	}
}

func TestRemoveEmptyDirFilesystem(t *testing.T) {
	tmp := t.TempDir()
	root := filepath.Join(tmp, "lib")
	store := storage.NewFilesystem(root)
	dir := filepath.Join("actors", "id1")
	_ = os.MkdirAll(filepath.Join(root, dir), 0o750)
	removeEmptyDirFilesystem(store, dir)
	if _, err := os.Stat(filepath.Join(root, dir)); !os.IsNotExist(err) {
		t.Fatal("expected directory removed")
	}
}
