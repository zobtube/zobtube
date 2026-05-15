package controller

import (
	"errors"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
)

// metadataStore returns storage for actor/channel/category thumbnails.
// migrated=false uses the default library; migrated=true uses metadata storage.
//
// Metadata migration worker: copy objects from metadataStore(false) to metadataStore(true),
// then set migrated=true on the row.
func (c *Controller) metadataStore(migrated bool) (storage.Storage, error) {
	if migrated {
		if c.metadataStorage == nil {
			return nil, errors.New("metadata storage not configured")
		}
		return c.metadataStorage, nil
	}
	if c.storageResolver == nil || c.config == nil {
		return nil, errors.New("library storage not configured")
	}
	return c.storageResolver.Storage(c.config.DefaultLibraryID)
}

// videoThumbnailStore returns storage for video/clip/movie thumbnails (not video.mp4).
// migrated=false uses the video's library; migrated=true uses metadata storage.
func (c *Controller) videoThumbnailStore(video *model.Video) (storage.Storage, error) {
	if video.Migrated {
		if c.metadataStorage == nil {
			return nil, errors.New("metadata storage not configured")
		}
		return c.metadataStorage, nil
	}
	if c.storageResolver == nil {
		return nil, errors.New("library storage not configured")
	}
	return c.storageResolver.Storage(c.videoLibraryID(video))
}

// metadataStoreForWrite returns metadata storage for new thumbnail uploads.
func (c *Controller) metadataStoreForWrite() (storage.Storage, error) {
	return c.metadataStore(true)
}
