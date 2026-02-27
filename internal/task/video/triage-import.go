package video

import (
	"errors"
	"io"
	"path/filepath"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

func importFromTriage(ctx *common.Context, params common.Parameters) (string, error) {
	id := params["videoID"]
	video := &model.Video{ID: id}
	result := ctx.DB.First(video)
	if result.RowsAffected < 1 {
		return "video does not exist", errors.New("id not in db")
	}
	libID := videoLibraryID(ctx, video)
	store, err := ctx.StorageResolver.Storage(libID)
	if err != nil {
		return "unable to resolve storage", err
	}
	triagePath := filepath.Join("triage", video.Filename)
	newPath := video.RelativePath()
	if err := store.MkdirAll(filepath.Dir(newPath)); err != nil {
		return "unable to create new video folder", err
	}
	rc, err := store.Open(triagePath)
	if err != nil {
		return "unable to open triage file", err
	}
	defer rc.Close()
	wc, err := store.Create(newPath)
	if err != nil {
		return "unable to create video file", err
	}
	defer wc.Close()
	if _, err := io.Copy(wc, rc); err != nil {
		return "unable to copy video", err
	}
	_ = store.Delete(triagePath)
	video.Imported = true
	if err := ctx.DB.Save(video).Error; err != nil {
		return "unable to update database", err
	}
	return "", nil
}

func videoLibraryID(ctx *common.Context, video *model.Video) string {
	if video.LibraryID != nil && *video.LibraryID != "" {
		return *video.LibraryID
	}
	return ctx.Config.DefaultLibraryID
}
