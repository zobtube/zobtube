package metamigrate

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
	"github.com/zobtube/zobtube/internal/task/common"
)

const TaskName = "metadata/migrate"

// NewMetadataMigrate returns a task that copies legacy thumbnails to metadata storage.
func NewMetadataMigrate() *common.Task {
	return &common.Task{
		Name: TaskName,
		Steps: []common.Step{
			{
				Name:     "migrate-all",
				NiceName: "Migrate thumbnails to metadata storage",
				Func:     migrateAll,
			},
		},
	}
}

type migrateStats struct {
	actors     int
	channels   int
	categories int
	videos     int
	skipped    int
}

func migrateAll(ctx *common.Context, _ common.Parameters) (string, error) {
	srcLib, err := legacyEntityStore(ctx)
	if err != nil {
		return "legacy library storage unavailable", err
	}
	dst, err := metadataTargetStore(ctx)
	if err != nil {
		return "metadata storage unavailable", err
	}

	var stats migrateStats

	if err := migrateActors(ctx, srcLib, dst, &stats); err != nil {
		return "actor migration failed", err
	}
	if err := migrateChannels(ctx, srcLib, dst, &stats); err != nil {
		return "channel migration failed", err
	}
	if err := migrateCategories(ctx, srcLib, dst, &stats); err != nil {
		return "category migration failed", err
	}
	cleanupLegacyEntityRoots(srcLib)
	if err := migrateVideos(ctx, dst, &stats); err != nil {
		return "video migration failed", err
	}

	return fmt.Sprintf(
		"migrated actors=%d channels=%d categories=%d videos=%d skipped=%d",
		stats.actors, stats.channels, stats.categories, stats.videos, stats.skipped,
	), nil
}

func migrateObject(src, dst storage.Storage, path string) error {
	exists, err := src.Exists(path)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("source object not found")
	}
	return copyStorageObject(src, dst, path)
}

func migrateActors(ctx *common.Context, src, dst storage.Storage, stats *migrateStats) error {
	var actors []model.Actor
	if err := ctx.DB.Where("thumbnail = ? AND migrated = ?", true, false).Find(&actors).Error; err != nil {
		return err
	}
	for _, actor := range actors {
		path := filepath.Join("actors", actor.ID, "thumb.jpg")
		if err := migrateObject(src, dst, path); err != nil {
			stats.skipped++
			continue
		}
		actor.Migrated = true
		if err := ctx.DB.Model(&actor).Update("migrated", true).Error; err != nil {
			return err
		}
		_ = removeStorageObject(src, path)
		removeEmptyDirFilesystem(src, filepath.Join("actors", actor.ID))
		stats.actors++
	}
	return nil
}

func migrateChannels(ctx *common.Context, src, dst storage.Storage, stats *migrateStats) error {
	var channels []model.Channel
	if err := ctx.DB.Where("thumbnail = ? AND migrated = ?", true, false).Find(&channels).Error; err != nil {
		return err
	}
	for _, ch := range channels {
		path := filepath.Join("channels", ch.ID, "thumb.jpg")
		if err := migrateObject(src, dst, path); err != nil {
			stats.skipped++
			continue
		}
		if err := ctx.DB.Model(&ch).Update("migrated", true).Error; err != nil {
			return err
		}
		_ = removeStorageObject(src, path)
		removeEmptyDirFilesystem(src, filepath.Join("channels", ch.ID))
		stats.channels++
	}
	return nil
}

func migrateCategories(ctx *common.Context, src, dst storage.Storage, stats *migrateStats) error {
	var subs []model.CategorySub
	if err := ctx.DB.Where("thumbnail = ? AND migrated = ?", true, false).Find(&subs).Error; err != nil {
		return err
	}
	for _, sub := range subs {
		path := filepath.Join("categories", fmt.Sprintf("%s.jpg", sub.ID))
		if err := migrateObject(src, dst, path); err != nil {
			stats.skipped++
			continue
		}
		if err := ctx.DB.Model(&sub).Update("migrated", true).Error; err != nil {
			return err
		}
		_ = removeStorageObject(src, path)
		stats.categories++
	}
	return nil
}

func migrateVideos(ctx *common.Context, dst storage.Storage, stats *migrateStats) error {
	var videos []model.Video
	if err := ctx.DB.Where("migrated = ? AND (thumbnail = ? OR thumbnail_mini = ?)", false, true, true).Find(&videos).Error; err != nil {
		return err
	}
	for _, video := range videos {
		src, err := legacyVideoStore(ctx, &video)
		if err != nil {
			stats.skipped++
			continue
		}
		migratedAny := false
		for _, path := range []string{video.ThumbnailRelativePath(), video.ThumbnailXSRelativePath()} {
			exists, err := src.Exists(path)
			if err != nil || !exists {
				continue
			}
			if err := copyStorageObject(src, dst, path); err != nil {
				return err
			}
			_ = removeStorageObject(src, path)
			migratedAny = true
		}
		if !migratedAny {
			stats.skipped++
			continue
		}
		if err := ctx.DB.Model(&video).Update("migrated", true).Error; err != nil {
			return err
		}
		cleanupVideoLegacyFolder(src, &video)
		stats.videos++
	}
	return nil
}

func cleanupVideoLegacyFolder(src storage.Storage, video *model.Video) {
	hasVideo, _ := src.Exists(video.RelativePath())
	if hasVideo {
		return
	}
	removeEmptyDirFilesystem(src, video.FolderRelativePath())
}
