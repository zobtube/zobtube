package video

import (
	"errors"
	"io"
	"path/filepath"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
	"github.com/zobtube/zobtube/internal/task/common"
)

func NewVideoMoveLibrary() *common.Task {
	return &common.Task{
		Name: "video/move-library",
		Steps: []common.Step{
			{Name: "move-files", NiceName: "Copy video and thumbnails to new library", Func: moveFilesToNewLibrary},
			{Name: "update-db", NiceName: "Update video library in database", Func: updateVideoLibraryID},
			{Name: "delete-source", NiceName: "Remove files from source library", Func: deleteFromSourceLibrary},
		},
	}
}

// moveFilesToNewLibrary copies video file and thumbnails from source storage to target storage.
func moveFilesToNewLibrary(ctx *common.Context, params common.Parameters) (string, error) {
	videoID := params["videoID"]
	targetLibID := params["targetLibraryID"]
	if targetLibID == "" {
		return "targetLibraryID required", errors.New("missing targetLibraryID")
	}
	video := &model.Video{ID: videoID}
	if ctx.DB.First(video).RowsAffected < 1 {
		return "video does not exist", errors.New("video not found")
	}
	sourceLibID := videoLibraryID(ctx, video)
	if sourceLibID == targetLibID {
		return "video already in target library", errors.New("same library")
	}
	sourceStore, err := ctx.StorageResolver.Storage(sourceLibID)
	if err != nil {
		return "unable to resolve source storage", err
	}
	targetStore, err := ctx.StorageResolver.Storage(targetLibID)
	if err != nil {
		return "unable to resolve target storage", err
	}
	videoPath, ok, err := storage.FirstExistingPath(sourceStore, video.StoragePathCandidates())
	if err != nil {
		return "unable to check source file", err
	}
	if !ok {
		return "video file not found on source library", errors.New("source video missing")
	}
	paths := []string{videoPath}
	if !video.Migrated {
		paths = append(paths, video.ThumbnailRelativePath(), video.ThumbnailXSRelativePath())
	}
	copied := 0
	for _, p := range paths {
		if p != videoPath {
			ok, err := sourceStore.Exists(p)
			if err != nil {
				return "unable to check source file", err
			}
			if !ok {
				continue
			}
		}
		if err := targetStore.MkdirAll(filepath.Dir(p)); err != nil {
			return "unable to create target folder", err
		}
		rc, err := sourceStore.Open(p)
		if err != nil {
			return "unable to open source file", err
		}
		wc, err := targetStore.Create(p)
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
		copied++
	}
	if copied == 0 {
		return "video file was not copied", errors.New("no video file copied")
	}
	params["sourceVideoPath"] = videoPath
	return "", nil
}

// updateVideoLibraryID sets video.LibraryID to the target library and saves.
func updateVideoLibraryID(ctx *common.Context, params common.Parameters) (string, error) {
	videoID := params["videoID"]
	targetLibID := params["targetLibraryID"]
	video := &model.Video{ID: videoID}
	if ctx.DB.First(video).RowsAffected < 1 {
		return "video does not exist", errors.New("video not found")
	}
	video.LibraryID = &targetLibID
	if err := ctx.DB.Save(video).Error; err != nil {
		return "unable to update video library", err
	}
	return "", nil
}

// deleteFromSourceLibrary removes the video file and thumbnails from the source library.
// Video record already points to target library at this step.
func deleteFromSourceLibrary(ctx *common.Context, params common.Parameters) (string, error) {
	videoID := params["videoID"]
	targetLibID := params["targetLibraryID"]
	video := &model.Video{ID: videoID}
	if ctx.DB.First(video).RowsAffected < 1 {
		return "video does not exist", errors.New("video not found")
	}
	// Source library is the one we moved from; video.LibraryID is now target, so we need
	// to get source from params. We don't have sourceLibraryID in params. So we need to
	// pass it in params from the first step - but params are fixed when task is created.
	// So we have two options: (1) pass sourceLibraryID in params when creating the task
	// (controller knows current library before creating task), or (2) in step 1 store
	// sourceLibraryID somewhere (e.g. in task.Parameters by updating the task record).
	// Option 1 is cleaner: add sourceLibraryID to the task params when creating the task.
	sourceLibID := params["sourceLibraryID"]
	if sourceLibID == "" {
		// Backward compat: if not set, skip delete (task was created before we added this param)
		return "", nil
	}
	if sourceLibID == targetLibID {
		return "", nil
	}
	store, err := ctx.StorageResolver.Storage(sourceLibID)
	if err != nil {
		return "unable to resolve source storage", err
	}
	sourceVideoPath := params["sourceVideoPath"]
	if sourceVideoPath == "" {
		sourceVideoPath, _, _ = storage.FirstExistingPath(store, video.StoragePathCandidates())
	}
	if sourceVideoPath != "" {
		_ = store.Delete(sourceVideoPath)
	}
	if !video.Migrated {
		_ = store.Delete(video.ThumbnailRelativePath())
		_ = store.Delete(video.ThumbnailXSRelativePath())
	}
	return "", nil
}
