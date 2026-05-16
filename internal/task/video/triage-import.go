package video

import (
	"errors"
	"io"
	"path/filepath"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

// importFromTriage moves an uploaded file from the library triage folder to
// the location dictated by the currently Active Organization, unless the
// import is requested with reorganization disabled (per-import flag or
// global setting). In the no-reorg case the file stays at
// triage/<filename> and Path is recorded so the rest of the pipeline can
// stream it from there.
func importFromTriage(ctx *common.Context, params common.Parameters) (string, error) {
	id := params["videoID"]
	video := &model.Video{ID: id}
	result := ctx.DB.First(video)
	if result.RowsAffected < 1 {
		return "video does not exist", errors.New("id not in db")
	}
	libID := videoLibraryID(ctx, video)
	store, err := ctx.StorageResolver.Storage(libID)
	if err != nil {
		return "unable to resolve storage", err
	}
	triagePath := filepath.Join("triage", video.Filename)

	skipReorg, err := shouldSkipReorganization(ctx, params)
	if err != nil {
		return "unable to resolve reorganization setting", err
	}

	if skipReorg {
		path := triagePath
		video.Path = &path
		video.OrganizationID = nil
		video.Imported = true
		if err := ctx.DB.Save(video).Error; err != nil {
			return "unable to update database", err
		}
		return "", nil
	}

	org, err := resolveActiveOrganization(ctx)
	if err != nil {
		return "unable to resolve active organization", err
	}
	newPath := org.Render(video)
	if err := store.MkdirAll(filepath.Dir(newPath)); err != nil {
		return "unable to create new video folder", err
	}
	rc, err := store.Open(triagePath)
	if err != nil {
		return "unable to open triage file", err
	}
	defer rc.Close()
	wc, err := store.Create(newPath)
	if err != nil {
		return "unable to create video file", err
	}
	defer wc.Close()
	if _, err := io.Copy(wc, rc); err != nil {
		return "unable to copy video", err
	}
	_ = store.Delete(triagePath)
	orgID := org.ID
	video.OrganizationID = &orgID
	video.Path = &newPath
	video.Imported = true
	if err := ctx.DB.Save(video).Error; err != nil {
		return "unable to update database", err
	}
	return "", nil
}

// shouldSkipReorganization returns true when the file should remain at its
// triage path. Per-task params win over the global Configuration default.
func shouldSkipReorganization(ctx *common.Context, params common.Parameters) (bool, error) {
	if v, ok := params["skipReorganization"]; ok && v != "" {
		return v == "true" || v == "1", nil
	}
	if ctx.DB == nil {
		return false, nil
	}
	var cfg model.Configuration
	if err := ctx.DB.First(&cfg).Error; err != nil {
		return false, nil
	}
	return !cfg.ReorganizeOnImport, nil
}

// resolveActiveOrganization returns the currently Active organization,
// falling back to the first one found if none are explicitly active.
func resolveActiveOrganization(ctx *common.Context) (*model.Organization, error) {
	var org model.Organization
	if err := ctx.DB.Where("active = ?", true).First(&org).Error; err == nil {
		return &org, nil
	}
	if err := ctx.DB.First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func videoLibraryID(ctx *common.Context, video *model.Video) string {
	if video.LibraryID != nil && *video.LibraryID != "" {
		return *video.LibraryID
	}
	return ctx.Config.DefaultLibraryID
}
