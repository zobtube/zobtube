package video

import (
	"errors"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

func NewVideoDeleting() *common.Task {
	return &common.Task{
		Name: "video/delete",
		Steps: []common.Step{
			{Name: "delete-thumbnail", NiceName: "Delete video's thumbnail", Func: deleteThumbnail},
			{Name: "delete-thumbnail-mini", NiceName: "Delete video's mini thumbnail", Func: deleteThumbnailMini},
			{Name: "delete-video", NiceName: "Delete video", Func: deleteVideoFile},
			{Name: "delete-video-folder", NiceName: "Delete video's folder", Func: deleteFolder},
			{Name: "delete-video-in-db", NiceName: "Delete video's database entry", Func: deleteInDatabase},
		},
	}
}

func deleteVideoFile(ctx *common.Context, params common.Parameters) (string, error) {
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
	_ = store.Delete(video.RelativePath())
	return "", nil
}

func deleteFolder(ctx *common.Context, params common.Parameters) (string, error) {
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
	// Delete known files in folder (storage has no recursive delete)
	_ = store.Delete(video.ThumbnailRelativePath())
	_ = store.Delete(video.ThumbnailXSRelativePath())
	_ = store.Delete(video.RelativePath())
	return "", nil
}

func deleteInDatabase(ctx *common.Context, params common.Parameters) (string, error) {
	videoID := params["videoID"]
	video := &model.Video{ID: videoID}
	result := ctx.DB.First(video)
	if result.RowsAffected < 1 {
		return "video does not exist", errors.New("id not in db")
	}
	if err := ctx.DB.Delete(&video).Error; err != nil {
		return "unable to delete video", err
	}
	return "", nil
}
