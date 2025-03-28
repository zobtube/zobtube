package video

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

func NewVideoDeleting() *common.Task {
	task := &common.Task{
		Name: "video/delete",
		Steps: []common.Step{
			{
				Name:     "delete-thumbnail",
				NiceName: "Delete video's thumbnail",
				Func:     deleteThumbnail,
			},
			{
				Name:     "delete-thumbnail-mini",
				NiceName: "Delete video's mini thumbnail",
				Func:     deleteThumbnailMini,
			},
			{
				Name:     "delete-video",
				NiceName: "Delete video",
				Func:     deleteVideoFile,
			},
			{
				Name:     "delete-video-folder",
				NiceName: "Delete video's folder",
				Func:     deleteFolder,
			},
			{
				Name:     "delete-video-in-db",
				NiceName: "Delete video's database entry",
				Func:     deleteInDatabase,
			},
		},
	}

	return task
}

func deleteVideoFile(ctx *common.Context, params common.Parameters) (string, error) {
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

	// check video presence
	videoPath := filepath.Join(ctx.Config.Media.Path, video.RelativePath())
	_, err := os.Stat(videoPath)
	if err != nil && !os.IsNotExist(err) {
		return "unable to check video presence on disk", err
	}
	if !os.IsNotExist(err) {
		// exist, deleting it
		err = os.Remove(videoPath)
		if err != nil {
			return "unable to delete video on disk", err
		}
	}

	return "", nil
}

func deleteFolder(ctx *common.Context, params common.Parameters) (string, error) {
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

	// delete folder
	folderPath := filepath.Join(ctx.Config.Media.Path, video.FolderRelativePath())
	_, err := os.Stat(folderPath)
	if err != nil && !os.IsNotExist(err) {
		return "unable to check video folder presence", err
	}
	if !os.IsNotExist(err) {
		// exist, deleting it
		err = os.Remove(folderPath)
		if err != nil {
			return "unable to delete video folder", err
		}
	}

	return "", nil
}

func deleteInDatabase(ctx *common.Context, params common.Parameters) (string, error) {
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

	err := ctx.DB.Delete(&video).Error
	if err != nil {
		return "unable to delete video", err
	}

	return "", nil
}
