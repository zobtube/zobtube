package video

import (
	"errors"
	"os/exec"
	"strings"
	"time"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
	"github.com/zobtube/zobtube/internal/task/common"
)

func computeDuration(ctx *common.Context, params common.Parameters) (string, error) {
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
	localPath, cleanup, err := storage.LocalPathForRead(store, video.RelativePath())
	if err != nil {
		return "unable to get local path for video", err
	}
	defer cleanup()
	// #nosec G204
	out, err := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		localPath,
	).Output()
	if err != nil {
		return "unable to retrieve video length", err
	}
	duration := strings.TrimSpace(string(out))
	d, err := time.ParseDuration(duration + "s")
	if err != nil {
		return "unable to parse duration", err
	}
	video.Duration = d
	if err := ctx.DB.Save(&video).Error; err != nil {
		return "unable to save video duration", err
	}
	return "", nil
}
