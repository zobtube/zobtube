package metamigrate

import (
	"errors"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
	"github.com/zobtube/zobtube/internal/task/common"
)

func legacyEntityStore(ctx *common.Context) (storage.Storage, error) {
	if ctx.StorageResolver == nil || ctx.Config == nil {
		return nil, errors.New("library storage not configured")
	}
	return ctx.StorageResolver.Storage(ctx.Config.DefaultLibraryID)
}

func metadataTargetStore(ctx *common.Context) (storage.Storage, error) {
	if ctx.MetadataStorage == nil {
		return nil, errors.New("metadata storage not configured")
	}
	return ctx.MetadataStorage, nil
}

func legacyVideoStore(ctx *common.Context, video *model.Video) (storage.Storage, error) {
	if ctx.StorageResolver == nil {
		return nil, errors.New("library storage not configured")
	}
	libID := videoLibraryID(ctx, video)
	return ctx.StorageResolver.Storage(libID)
}

func videoLibraryID(ctx *common.Context, video *model.Video) string {
	if video.LibraryID != nil && *video.LibraryID != "" {
		return *video.LibraryID
	}
	return ctx.Config.DefaultLibraryID
}
