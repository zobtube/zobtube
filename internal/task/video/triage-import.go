package video

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

func importFromTriage(ctx *common.Context, params common.Parameters) (string, error) {
	// get id from path
	id := params["videoID"]

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := ctx.DB.First(video)

	// check result
	if result.RowsAffected < 1 {
		return "video does not exist", errors.New("id not in db")
	}

	// prepare paths
	previousPath := filepath.Join(ctx.Config.Media.Path, "/triage", video.Filename)
	newFolderPath := filepath.Join(ctx.Config.Media.Path, video.FolderRelativePath())
	newPath := filepath.Join(ctx.Config.Media.Path, video.RelativePath())

	// ensure folder exists
	_, err := os.Stat(newFolderPath)
	if os.IsNotExist(err) {
		// do not exists, create it
		err = os.Mkdir(newFolderPath, os.ModePerm)
		if err != nil {
			return "unable to create new video folder", err
		}
	} else if err != nil {
		return "unable to read new video folder", err
	}

	// move
	err = os.Rename(previousPath, newPath)
	if err != nil {
		return "unable to move new video into its folder", err
	}

	// commit the update on database
	video.Imported = true
	err = ctx.DB.Save(video).Error
	if err != nil {
		return "unable to update database", err
	}

	return "", nil
}
