package video

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
	"github.com/zobtube/zobtube/internal/task/common"
)

func TestReorganize_EndToEnd_MovesFileAndUpdatesRow(t *testing.T) {
	tmp := t.TempDir()
	ctx := setupTaskContext(t, tmp)

	// Default org from bootstrap matches the legacy layout.
	srcOrg := &model.Organization{}
	if err := ctx.DB.First(srcOrg, "id = ?", model.DefaultOrganizationUUID).Error; err != nil {
		t.Fatal(err)
	}
	// A second organization to migrate to.
	dstOrg := model.Organization{
		ID:       "22222222-2222-2222-2222-222222222222",
		Name:     "v2 flat",
		Template: "media/$ID$EXT",
		Active:   false,
	}
	if err := ctx.DB.Create(&dstOrg).Error; err != nil {
		t.Fatal(err)
	}

	srcOrgID := srcOrg.ID
	srcPath := filepath.Join("videos", "vid1", "video.mp4")
	v := model.Video{
		ID:             "vid1",
		Name:           "v",
		Filename:       "raw.mp4",
		Type:           "v",
		Imported:       true,
		OrganizationID: &srcOrgID,
		Path:           &srcPath,
	}
	if err := ctx.DB.Create(&v).Error; err != nil {
		t.Fatal(err)
	}

	// Create the legacy file on the storage.
	libStore := storage.NewFilesystem(tmp)
	_ = libStore.MkdirAll(filepath.Dir(srcPath))
	wc, err := libStore.Create(srcPath)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = wc.Write([]byte("payload"))
	_ = wc.Close()

	params := common.Parameters{
		"videoID":              v.ID,
		"targetOrganizationID": dstOrg.ID,
		"sourcePath":           srcPath,
	}

	if msg, err := reorganizeMoveFile(ctx, params); err != nil {
		t.Fatalf("move-file: %v (%s)", err, msg)
	}
	expectedDst := "media/vid1.mp4"
	if _, err := os.Stat(filepath.Join(tmp, expectedDst)); err != nil {
		t.Fatalf("expected destination file at %s: %v", expectedDst, err)
	}

	if msg, err := reorganizeUpdateDB(ctx, params); err != nil {
		t.Fatalf("update-db: %v (%s)", err, msg)
	}
	var afterUpdate model.Video
	if err := ctx.DB.First(&afterUpdate, "id = ?", v.ID).Error; err != nil {
		t.Fatal(err)
	}
	if afterUpdate.OrganizationID == nil || *afterUpdate.OrganizationID != dstOrg.ID {
		t.Errorf("expected organization_id %q, got %v", dstOrg.ID, afterUpdate.OrganizationID)
	}
	if afterUpdate.Path == nil || *afterUpdate.Path != expectedDst {
		t.Errorf("expected path %q, got %v", expectedDst, afterUpdate.Path)
	}

	if msg, err := reorganizeDeleteSource(ctx, params); err != nil {
		t.Fatalf("delete-source: %v (%s)", err, msg)
	}
	if _, err := os.Stat(filepath.Join(tmp, srcPath)); !os.IsNotExist(err) {
		t.Errorf("expected source file removed, stat err = %v", err)
	}
}

func TestReorganize_NoopWhenSourceMatchesTarget(t *testing.T) {
	tmp := t.TempDir()
	ctx := setupTaskContext(t, tmp)
	srcOrg := &model.Organization{}
	_ = ctx.DB.First(srcOrg, "id = ?", model.DefaultOrganizationUUID).Error

	srcOrgID := srcOrg.ID
	path := filepath.Join("videos", "vid1", "video.mp4")
	v := model.Video{
		ID:             "vid1",
		Name:           "v",
		Filename:       "raw.mp4",
		Type:           "v",
		Imported:       true,
		OrganizationID: &srcOrgID,
		Path:           &path,
	}
	if err := ctx.DB.Create(&v).Error; err != nil {
		t.Fatal(err)
	}

	// Reorganizing to the same org should be a no-op.
	params := common.Parameters{
		"videoID":              v.ID,
		"targetOrganizationID": srcOrg.ID,
	}
	if msg, err := reorganizeMoveFile(ctx, params); err != nil {
		t.Fatalf("expected no-op, got %v (%s)", err, msg)
	}
}
