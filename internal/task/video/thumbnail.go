package video

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
	"github.com/zobtube/zobtube/internal/task/common"
)

func generateThumbnail(ctx *common.Context, params common.Parameters) (string, error) {
	id := params["videoID"]
	timing := params["thumbnailTiming"]
	video := &model.Video{ID: id}
	result := ctx.DB.First(video)
	if result.RowsAffected < 1 {
		return "video does not exist", errors.New("id not in db")
	}
	store, err := ctx.StorageResolver.Storage(videoLibraryID(ctx, video))
	if err != nil {
		return "unable to resolve storage", err
	}
	videoLocal, cleanupVideo, err := storage.LocalPathForRead(store, video.RelativePath())
	if err != nil {
		return "unable to get local path for video", err
	}
	defer cleanupVideo()
	thumbPath := video.ThumbnailRelativePath()
	thumbTemp, err := os.CreateTemp("", "zt-thumb-*.jpg")
	if err != nil {
		return "unable to create temp for thumbnail", err
	}
	thumbTempPath := thumbTemp.Name()
	thumbTemp.Close()
	defer os.Remove(thumbTempPath)
	// #nosec G204
	_, err = exec.Command(
		"ffmpeg", "-y", "-ss", timing,
		"-i", videoLocal,
		"-frames:v", "1", "-q:v", "2",
		thumbTempPath,
	).Output()
	if err != nil {
		return "unable to generate thumbnail with ffmpeg", err
	}
	if err := store.MkdirAll(filepath.Dir(thumbPath)); err != nil {
		return "unable to create thumbnail folder", err
	}
	r, err := os.Open(thumbTempPath)
	if err != nil {
		return "unable to open temp thumbnail", err
	}
	defer r.Close()
	w, err := store.Create(thumbPath)
	if err != nil {
		return "unable to create thumbnail file", err
	}
	defer w.Close()
	if _, err := io.Copy(w, r); err != nil {
		return "unable to write thumbnail", err
	}
	video.Thumbnail = true
	ctx.DB.Save(&video)
	return "", nil
}

func deleteThumbnail(ctx *common.Context, params common.Parameters) (string, error) {
	videoID := params["videoID"]
	video := &model.Video{ID: videoID}
	result := ctx.DB.First(video)
	if result.RowsAffected < 1 {
		return "video does not exist", errors.New("id not in db")
	}
	store, err := ctx.StorageResolver.Storage(videoLibraryID(ctx, video))
	if err != nil {
		return "unable to resolve storage", err
	}
	thumbPath := video.ThumbnailRelativePath()
	_ = store.Delete(thumbPath)
	return "", nil
}
