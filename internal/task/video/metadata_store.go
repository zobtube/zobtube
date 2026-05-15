package video

import (
	"errors"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
	"github.com/zobtube/zobtube/internal/task/common"
)

func metadataStore(ctx *common.Context, migrated bool) (storage.Storage, error) {
	if migrated {
		if ctx.MetadataStorage == nil {
			return nil, errors.New("metadata storage not configured")
		}
		return ctx.MetadataStorage, nil
	}
	if ctx.StorageResolver == nil || ctx.Config == nil {
		return nil, errors.New("library storage not configured")
	}
	return ctx.StorageResolver.Storage(ctx.Config.DefaultLibraryID)
}

func videoThumbnailStore(ctx *common.Context, video *model.Video) (storage.Storage, error) {
	if video.Migrated {
		if ctx.MetadataStorage == nil {
			return nil, errors.New("metadata storage not configured")
		}
		return ctx.MetadataStorage, nil
	}
	if ctx.StorageResolver == nil {
		return nil, errors.New("library storage not configured")
	}
	return ctx.StorageResolver.Storage(videoLibraryID(ctx, video))
}

func metadataStoreForWrite(ctx *common.Context) (storage.Storage, error) {
	return metadataStore(ctx, true)
}
