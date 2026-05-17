package photoset

import (
	"errors"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

func NewPhotosetDeleting() *common.Task {
	return &common.Task{
		Name: "photoset/delete",
		Steps: []common.Step{
			{Name: "delete-photo-files", NiceName: "Delete photo files", Func: deletePhotoFiles},
			{Name: "delete-photo-thumbnails", NiceName: "Delete photo thumbnails", Func: deletePhotoThumbnails},
			{Name: "delete-photoset-folder", NiceName: "Delete photoset folder", Func: deletePhotosetFolder},
			{Name: "delete-photoset-in-db", NiceName: "Delete photoset database entry", Func: deletePhotosetInDB},
		},
	}
}

func deletePhotoFiles(ctx *common.Context, params common.Parameters) (string, error) {
	psID := params["photosetID"]
	ps := &model.Photoset{ID: psID}
	if ctx.DB.Preload("Photos").First(ps).RowsAffected < 1 {
		return "photoset does not exist", errors.New("not found")
	}
	store, err := ctx.StorageResolver.Storage(photosetLibraryID(ctx, ps))
	if err != nil {
		return "unable to resolve storage", err
	}
	for i := range ps.Photos {
		_ = store.Delete(ps.Photos[i].RelativePath(ps))
	}
	return "", nil
}

func deletePhotoThumbnails(ctx *common.Context, params common.Parameters) (string, error) {
	psID := params["photosetID"]
	ps := &model.Photoset{ID: psID}
	if ctx.DB.Preload("Photos").First(ps).RowsAffected < 1 {
		return "photoset does not exist", errors.New("not found")
	}
	meta, err := metadataStoreForWrite(ctx)
	if err != nil {
		return "unable to resolve metadata storage", err
	}
	for i := range ps.Photos {
		_ = meta.Delete(ps.Photos[i].ThumbnailMiniRelativePath(ps))
	}
	return "", nil
}

func deletePhotosetFolder(ctx *common.Context, params common.Parameters) (string, error) {
	psID := params["photosetID"]
	ps := &model.Photoset{ID: psID}
	if ctx.DB.First(ps).RowsAffected < 1 {
		return "photoset does not exist", errors.New("not found")
	}
	store, err := ctx.StorageResolver.Storage(photosetLibraryID(ctx, ps))
	if err != nil {
		return "unable to resolve storage", err
	}
	_ = store.Delete(ps.FolderRelativePath())
	return "", nil
}

func deletePhotosetInDB(ctx *common.Context, params common.Parameters) (string, error) {
	psID := params["photosetID"]
	ps := &model.Photoset{ID: psID}
	if ctx.DB.First(ps).RowsAffected < 1 {
		return "photoset does not exist", errors.New("not found")
	}
	if err := ctx.DB.Where("photoset_id = ?", psID).Delete(&model.Photo{}).Error; err != nil {
		return "unable to delete photos", err
	}
	if err := ctx.DB.Delete(ps).Error; err != nil {
		return "unable to delete photoset", err
	}
	return "", nil
}
