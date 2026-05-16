package video

import (
	"errors"
	"io"
	"path/filepath"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

// NewVideoReorganize returns the task that moves a single video file to the
// path produced by the target organization. The thumbnail layout is owned by
// metadata storage and is intentionally left untouched.
func NewVideoReorganize() *common.Task {
	return &common.Task{
		Name: "video/reorganize",
		Steps: []common.Step{
			{Name: "move-file", NiceName: "Copy video file to the target organization layout", Func: reorganizeMoveFile},
			{Name: "update-db", NiceName: "Update video path and organization in database", Func: reorganizeUpdateDB},
			{Name: "delete-source", NiceName: "Remove the previous video file", Func: reorganizeDeleteSource},
		},
	}
}

// reorganizeMoveFile copies the existing video file (located at the current
// Path or legacy hardcoded layout) to the location rendered from the target
// organization. The source file is left in place at this step so the move
// can be safely resumed if any step fails.
func reorganizeMoveFile(ctx *common.Context, params common.Parameters) (string, error) {
	videoID := params["videoID"]
	orgID := params["targetOrganizationID"]
	if orgID == "" {
		return "targetOrganizationID required", errors.New("missing targetOrganizationID")
	}
	video := &model.Video{ID: videoID}
	if ctx.DB.First(video).RowsAffected < 1 {
		return "video does not exist", errors.New("video not found")
	}
	if !video.Imported {
		return "video has not been imported yet, nothing to reorganize", errors.New("not imported")
	}
	org := &model.Organization{}
	if err := ctx.DB.First(org, "id = ?", orgID).Error; err != nil {
		return "target organization not found", err
	}
	srcPath := video.RelativePath()
	dstPath := org.Render(video)
	if srcPath == dstPath {
		return "", nil
	}
	store, err := ctx.StorageResolver.Storage(videoLibraryID(ctx, video))
	if err != nil {
		return "unable to resolve storage", err
	}
	exists, err := store.Exists(srcPath)
	if err != nil {
		return "unable to stat source file", err
	}
	if !exists {
		return "source file missing on storage", errors.New("source missing")
	}
	if err := store.MkdirAll(filepath.Dir(dstPath)); err != nil {
		return "unable to create target folder", err
	}
	rc, err := store.Open(srcPath)
	if err != nil {
		return "unable to open source file", err
	}
	wc, err := store.Create(dstPath)
	if err != nil {
		rc.Close()
		return "unable to create target file", err
	}
	_, err = io.Copy(wc, rc)
	rc.Close()
	if err != nil {
		wc.Close()
		return "unable to copy file", err
	}
	if err := wc.Close(); err != nil {
		return "unable to close target file", err
	}
	return "", nil
}

// reorganizeUpdateDB flips the row to point at the new path and target
// organization. After this step a stream/thumb request will hit the new
// location.
func reorganizeUpdateDB(ctx *common.Context, params common.Parameters) (string, error) {
	videoID := params["videoID"]
	orgID := params["targetOrganizationID"]
	video := &model.Video{ID: videoID}
	if ctx.DB.First(video).RowsAffected < 1 {
		return "video does not exist", errors.New("video not found")
	}
	org := &model.Organization{}
	if err := ctx.DB.First(org, "id = ?", orgID).Error; err != nil {
		return "target organization not found", err
	}
	newPath := org.Render(video)
	updates := map[string]any{
		"organization_id": orgID,
		"path":            newPath,
	}
	if err := ctx.DB.Model(video).Updates(updates).Error; err != nil {
		return "unable to update video", err
	}
	return "", nil
}

// reorganizeDeleteSource removes the previous file when the move succeeded.
// The original path is recomputed from "sourcePath" if provided, otherwise
// we trust the DB row that no longer points at the legacy location.
func reorganizeDeleteSource(ctx *common.Context, params common.Parameters) (string, error) {
	videoID := params["videoID"]
	sourcePath := params["sourcePath"]
	video := &model.Video{ID: videoID}
	if ctx.DB.First(video).RowsAffected < 1 {
		return "video does not exist", errors.New("video not found")
	}
	if sourcePath == "" {
		return "", nil
	}
	if sourcePath == video.RelativePath() {
		return "", nil
	}
	store, err := ctx.StorageResolver.Storage(videoLibraryID(ctx, video))
	if err != nil {
		return "unable to resolve storage", err
	}
	_ = store.Delete(sourcePath)
	return "", nil
}
