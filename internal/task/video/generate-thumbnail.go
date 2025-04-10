package video

import "github.com/zobtube/zobtube/internal/task/common"

func NewVideoGenerateThumbnail() *common.Task {
	task := &common.Task{
		Name: "video/generate-thumbnail",
		Steps: []common.Step{
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
		},
	}

	return task
}
