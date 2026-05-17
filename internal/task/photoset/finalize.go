package photoset

import (
	"errors"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

func NewPhotosetFinalize() *common.Task {
	return &common.Task{
		Name: "photoset/finalize",
		Steps: []common.Step{
			{Name: "index-photos", NiceName: "Index photos on disk", Func: indexPhotos},
			{Name: "assign-positions", NiceName: "Assign photo positions", Func: assignPositions},
			{Name: "generate-mini-thumbnails", NiceName: "Generate mini thumbnails", Func: generateMiniThumbnails},
			{Name: "assign-cover", NiceName: "Assign cover photo", Func: assignCover},
			{Name: "apply-organization", NiceName: "Apply organization layout", Func: applyOrganization},
			{Name: "creating-to-ready", NiceName: "Finalize photoset status", Func: creatingToReady},
		},
	}
}

func indexPhotos(ctx *common.Context, params common.Parameters) (string, error) {
	psID := params["photosetID"]
	ps := &model.Photoset{ID: psID}
	if ctx.DB.First(ps).RowsAffected < 1 {
		return "photoset does not exist", errors.New("not found")
	}
	store, err := ctx.StorageResolver.Storage(photosetLibraryID(ctx, ps))
	if err != nil {
		return "unable to resolve storage", err
	}
	folder := ps.FolderRelativePath()
	entries, err := store.List(folder)
	if err != nil {
		return "unable to list photoset folder", err
	}
	for _, entry := range entries {
		if entry.IsDir {
			continue
		}
		name := entry.Name
		if strings.HasPrefix(name, "_") {
			continue
		}
		base := sanitizeBasename(name)
		if base == "" {
			continue
		}
		mimeType, ok := photoMimeForFilename(base)
		if !ok {
			continue
		}
		rel := filepath.Join(folder, base)
		var existing model.Photo
		if ctx.DB.Where("photoset_id = ? AND filename = ?", ps.ID, base).First(&existing).RowsAffected > 0 {
			continue
		}
		path := rel
		photo := &model.Photo{
			PhotosetID: ps.ID,
			Filename:   base,
			Path:       &path,
			Mime:       mimeType,
		}
		if err := ctx.DB.Create(photo).Error; err != nil {
			return "unable to create photo row", err
		}
	}
	ps.Imported = true
	if err := ctx.DB.Save(ps).Error; err != nil {
		return "unable to update photoset", err
	}
	return "", nil
}

func assignPositions(ctx *common.Context, params common.Parameters) (string, error) {
	psID := params["photosetID"]
	var photos []model.Photo
	if err := ctx.DB.Where("photoset_id = ?", psID).Find(&photos).Error; err != nil {
		return "unable to load photos", err
	}
	sort.Slice(photos, func(i, j int) bool {
		return strings.ToLower(photos[i].Filename) < strings.ToLower(photos[j].Filename)
	})
	for i := range photos {
		photos[i].Position = i + 1
		if err := ctx.DB.Model(&photos[i]).Update("position", photos[i].Position).Error; err != nil {
			return "unable to update position", err
		}
	}
	return "", nil
}

func generateMiniThumbnails(ctx *common.Context, params common.Parameters) (string, error) {
	psID := params["photosetID"]
	ps := &model.Photoset{ID: psID}
	if ctx.DB.First(ps).RowsAffected < 1 {
		return "photoset does not exist", errors.New("not found")
	}
	var photos []model.Photo
	if err := ctx.DB.Where("photoset_id = ?", psID).Order("position asc, filename asc").Find(&photos).Error; err != nil {
		return "unable to load photos", err
	}
	for i := range photos {
		if photos[i].ThumbnailMini {
			continue
		}
		if msg, err := generatePhotoMiniThumbnail(ctx, ps, &photos[i]); err != nil {
			return msg, err
		}
	}
	return "", nil
}

func assignCover(ctx *common.Context, params common.Parameters) (string, error) {
	psID := params["photosetID"]
	ps := &model.Photoset{ID: psID}
	if ctx.DB.First(ps).RowsAffected < 1 {
		return "photoset does not exist", errors.New("not found")
	}
	if ps.CoverPhotoID != nil && *ps.CoverPhotoID != "" {
		return "", nil
	}
	var photo model.Photo
	if ctx.DB.Where("photoset_id = ?", psID).Order("position asc, filename asc").First(&photo).RowsAffected < 1 {
		return "", nil
	}
	ps.CoverPhotoID = &photo.ID
	if err := ctx.DB.Save(ps).Error; err != nil {
		return "unable to set cover", err
	}
	return "", nil
}

func applyOrganization(ctx *common.Context, params common.Parameters) (string, error) {
	psID := params["photosetID"]
	ps := &model.Photoset{ID: psID}
	if ctx.DB.Preload("Photos").First(ps).RowsAffected < 1 {
		return "photoset does not exist", errors.New("not found")
	}
	org, err := model.ActiveOrganizationForScope(ctx.DB, model.OrganizationScopePhotoset)
	if err != nil {
		return "no active photoset organization", err
	}
	orgID := org.ID
	ps.OrganizationID = &orgID
	store, err := ctx.StorageResolver.Storage(photosetLibraryID(ctx, ps))
	if err != nil {
		return "unable to resolve storage", err
	}
	for i := range ps.Photos {
		photo := &ps.Photos[i]
		srcPath := photo.RelativePath(ps)
		dstPath := org.RenderPhotoset(ps, photo)
		if srcPath == dstPath {
			path := dstPath
			photo.Path = &path
			_ = ctx.DB.Save(photo).Error
			continue
		}
		exists, err := store.Exists(srcPath)
		if err != nil || !exists {
			continue
		}
		if err := store.MkdirAll(filepath.Dir(dstPath)); err != nil {
			return "unable to create target folder", err
		}
		rc, err := store.Open(srcPath)
		if err != nil {
			return "unable to open source file", err
		}
		wc, err := store.Create(dstPath)
		if err != nil {
			rc.Close()
			return "unable to create target file", err
		}
		_, err = io.Copy(wc, rc)
		rc.Close()
		if err != nil {
			wc.Close()
			return "unable to copy file", err
		}
		if err := wc.Close(); err != nil {
			return "unable to close target file", err
		}
		if srcPath != dstPath {
			_ = store.Delete(srcPath)
		}
		path := dstPath
		photo.Path = &path
		if err := ctx.DB.Save(photo).Error; err != nil {
			return "unable to update photo path", err
		}
	}
	if err := ctx.DB.Model(ps).Update("organization_id", orgID).Error; err != nil {
		return "unable to update photoset organization", err
	}
	return "", nil
}

func creatingToReady(ctx *common.Context, params common.Parameters) (string, error) {
	psID := params["photosetID"]
	ps := &model.Photoset{ID: psID}
	if ctx.DB.First(ps).RowsAffected < 1 {
		return "photoset does not exist", errors.New("not found")
	}
	ps.Status = model.PhotosetStatusReady
	if err := ctx.DB.Save(ps).Error; err != nil {
		return "unable to update status", err
	}
	return "", nil
}
