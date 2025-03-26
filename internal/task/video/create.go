package video

import (
	"errors"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

func NewVideoCreating() *common.Task {
	task := &common.Task{
		Name: "video/create",
		Steps: []common.Step{
			{
				Name:     "triage-import",
				NiceName: "Move the video from triage to its own folder",
				Func:     importFromTriage,
			},
			{
				Name:     "compute-duration",
				NiceName: "Compute duration of the video",
				Func:     computeDuration,
			},
			{
				Name:     "generate-thumbnail",
				NiceName: "Generate thumbnail",
				Func:     generateThumbnail,
			},
			{
				Name:     "generate-thumbnail-mini",
				NiceName: "Generate mini thumbnail",
				Func:     generateThumbnailMini,
			},
			{
				Name:     "creating-to-ready",
				NiceName: "Finalize DB status",
				Func:     creatingToReady,
			},
		},
	}

	return task
}

func creatingToReady(ctx *common.Context, params common.Parameters) (string, error) {
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

	video.Status = model.VideoStatusCreating

	err := ctx.DB.Save(&video).Error
	if err != nil {
		return "unable to save video duration", err
	}

	return "", nil
}
