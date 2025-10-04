package video

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

func computeDuration(ctx *common.Context, params common.Parameters) (string, error) {
	videoID := params["videoID"]

	// get item from ID
	video := &model.Video{
		ID: videoID,
	}
	result := ctx.DB.First(video)

	// check result
	if result.RowsAffected < 1 {
		return "video does not exist", errors.New("id not in db")
	}

	filePath := filepath.Join(ctx.Config.Media.Path, video.RelativePath())
	// #nosec G204
	out, err := exec.Command(
		"ffprobe",
		"-v",
		"error",
		"-show_entries",
		"format=duration",
		"-of",
		"default=noprint_wrappers=1:nokey=1",
		filePath,
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
	err = ctx.DB.Save(&video).Error
	if err != nil {
		return "unable to save video duration", err
	}

	return "", nil
}
