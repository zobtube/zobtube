package video

import (
	"path/filepath"
	"testing"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
	"github.com/zobtube/zobtube/internal/task/common"
)

// Regression: in-place import keeps file at triage/<filename>; move-library must copy
// from that path, not from the legacy videos/<id>/video.mp4 layout.
func TestMoveFilesToNewLibrary_InPlaceTriagePath(t *testing.T) {
	srcRoot := t.TempDir()
	dstRoot := t.TempDir()
	ctx := setupTaskContext(t, srcRoot)

	srcLib := model.Library{
		Name: "Source",
		Type: model.LibraryTypeFilesystem,
		Config: model.LibraryConfig{
			Filesystem: &model.LibraryConfigFilesystem{Path: srcRoot},
		},
	}
	dstLib := model.Library{
		Name: "Target",
		Type: model.LibraryTypeFilesystem,
		Config: model.LibraryConfig{
			Filesystem: &model.LibraryConfigFilesystem{Path: dstRoot},
		},
	}
	if err := ctx.DB.Create(&srcLib).Error; err != nil {
		t.Fatal(err)
	}
	if err := ctx.DB.Create(&dstLib).Error; err != nil {
		t.Fatal(err)
	}

	triagePath := filepath.Join("triage", "clip.mp4")
	srcStore := storage.NewFilesystem(srcRoot)
	_ = srcStore.MkdirAll("triage")
	wc, err := srcStore.Create(triagePath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := wc.Write([]byte("video-bytes")); err != nil {
		t.Fatal(err)
	}
	_ = wc.Close()

	path := triagePath
	v := &model.Video{
		ID:        "22222222-2222-2222-2222-222222222222",
		Name:      "In place",
		Filename:  "clip.mp4",
		Type:      "v",
		Imported:  true,
		Path:      &path,
		LibraryID: &srcLib.ID,
	}
	if err := ctx.DB.Create(v).Error; err != nil {
		t.Fatal(err)
	}

	if msg, err := moveFilesToNewLibrary(ctx, common.Parameters{
		"videoID":         v.ID,
		"targetLibraryID": dstLib.ID,
	}); err != nil {
		t.Fatalf("moveFilesToNewLibrary: %v (%s)", err, msg)
	}

	dstStore := storage.NewFilesystem(dstRoot)
	ok, err := dstStore.Exists(triagePath)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("expected file at %q on target library, legacy path would be %q", triagePath, v.RelativePath())
	}
}

// When Path is unset but the file is in triage/, move-library must still find and copy it.
func TestMoveFilesToNewLibrary_FindsTriageWhenDBPathWrong(t *testing.T) {
	srcRoot := t.TempDir()
	dstRoot := t.TempDir()
	ctx := setupTaskContext(t, srcRoot)

	srcLib := model.Library{
		Name: "Source",
		Type: model.LibraryTypeFilesystem,
		Config: model.LibraryConfig{
			Filesystem: &model.LibraryConfigFilesystem{Path: srcRoot},
		},
	}
	dstLib := model.Library{
		Name: "Target",
		Type: model.LibraryTypeFilesystem,
		Config: model.LibraryConfig{
			Filesystem: &model.LibraryConfigFilesystem{Path: dstRoot},
		},
	}
	if err := ctx.DB.Create(&srcLib).Error; err != nil {
		t.Fatal(err)
	}
	if err := ctx.DB.Create(&dstLib).Error; err != nil {
		t.Fatal(err)
	}

	triagePath := filepath.Join("triage", "orphan.mp4")
	srcStore := storage.NewFilesystem(srcRoot)
	_ = srcStore.MkdirAll("triage")
	wc, _ := srcStore.Create(triagePath)
	_, _ = wc.Write([]byte("video-bytes"))
	_ = wc.Close()

	v := &model.Video{
		ID:        "33333333-3333-3333-3333-333333333333",
		Name:      "Orphan",
		Filename:  "orphan.mp4",
		Type:      "v",
		Imported:  true,
		LibraryID: &srcLib.ID,
	}
	if err := ctx.DB.Create(v).Error; err != nil {
		t.Fatal(err)
	}

	if msg, err := moveFilesToNewLibrary(ctx, common.Parameters{
		"videoID":         v.ID,
		"targetLibraryID": dstLib.ID,
	}); err != nil {
		t.Fatalf("moveFilesToNewLibrary: %v (%s)", err, msg)
	}

	dstStore := storage.NewFilesystem(dstRoot)
	ok, err := dstStore.Exists(triagePath)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("expected file copied to target at %q", triagePath)
	}
}

func TestMoveFilesToNewLibrary_FailsWhenNoSourceFile(t *testing.T) {
	srcRoot := t.TempDir()
	dstRoot := t.TempDir()
	ctx := setupTaskContext(t, srcRoot)

	srcLib := model.Library{
		Name: "Source",
		Type: model.LibraryTypeFilesystem,
		Config: model.LibraryConfig{
			Filesystem: &model.LibraryConfigFilesystem{Path: srcRoot},
		},
	}
	dstLib := model.Library{
		Name: "Target",
		Type: model.LibraryTypeFilesystem,
		Config: model.LibraryConfig{
			Filesystem: &model.LibraryConfigFilesystem{Path: dstRoot},
		},
	}
	if err := ctx.DB.Create(&srcLib).Error; err != nil {
		t.Fatal(err)
	}
	if err := ctx.DB.Create(&dstLib).Error; err != nil {
		t.Fatal(err)
	}

	v := &model.Video{
		ID:        "44444444-4444-4444-4444-444444444444",
		Name:      "Missing",
		Filename:  "missing.mp4",
		Type:      "v",
		Imported:  true,
		LibraryID: &srcLib.ID,
	}
	if err := ctx.DB.Create(v).Error; err != nil {
		t.Fatal(err)
	}

	_, err := moveFilesToNewLibrary(ctx, common.Parameters{
		"videoID":         v.ID,
		"targetLibraryID": dstLib.ID,
	})
	if err == nil {
		t.Fatal("expected error when source video is missing")
	}
}
