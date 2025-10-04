package video

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

func generateThumbnail(ctx *common.Context, params common.Parameters) (string, error) {
	// get id from path
	id := params["videoID"]

	// get timing from path
	timing := params["thumbnailTiming"]

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := ctx.DB.First(video)

	// check result
	if result.RowsAffected < 1 {
		return "video does not exist", errors.New("id not in db")
	}

	// construct paths
	videoPath := filepath.Join(ctx.Config.Media.Path, video.RelativePath())
	videoPath, err := filepath.Abs(videoPath)
	if err != nil {
		return "Unable to get absolute path of video", err
	}

	thumbPath := filepath.Join(ctx.Config.Media.Path, video.ThumbnailRelativePath())
	thumbPath, err = filepath.Abs(thumbPath)
	if err != nil {
		return "Unable to get absolute path of the new thumbnail", err
	}

	// #nosec G204
	_, err = exec.Command(
		"ffmpeg",
		"-y",
		"-ss",
		timing,
		"-i",
		videoPath,
		"-frames:v",
		"1",
		"-q:v",
		"2",
		thumbPath,
	).Output()
	if err != nil {
		return "Unable to generate thumbnail with ffmpeg", err
	}

	video.Thumbnail = true
	ctx.DB.Save(&video)

	return "", nil
}

func deleteThumbnail(ctx *common.Context, params common.Parameters) (string, error) {
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

	// check thumb presence
	thumbPath := filepath.Join(ctx.Config.Media.Path, video.ThumbnailRelativePath())
	_, err := os.Stat(thumbPath)
	if err != nil && !os.IsNotExist(err) {
		return "unable to check thumbnail presence", err
	}
	if !os.IsNotExist(err) {
		// exist, deleting it
		err = os.Remove(thumbPath)
		if err != nil {
			return "unable to delete thumbnail", err
		}
	}

	return "", nil
}
