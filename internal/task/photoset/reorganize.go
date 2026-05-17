package photoset

import (
	"errors"
	"io"
	"path/filepath"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
)

func NewPhotosetReorganize() *common.Task {
	return &common.Task{
		Name: "photoset/reorganize",
		Steps: []common.Step{
			{Name: "move-photos", NiceName: "Move photos to target organization layout", Func: reorganizeMovePhotos},
			{Name: "update-db", NiceName: "Update photoset organization in database", Func: reorganizeUpdateDB},
		},
	}
}

func reorganizeMovePhotos(ctx *common.Context, params common.Parameters) (string, error) {
	psID := params["photosetID"]
	orgID := params["targetOrganizationID"]
	if orgID == "" {
		return "targetOrganizationID required", errors.New("missing targetOrganizationID")
	}
	ps := &model.Photoset{ID: psID}
	if ctx.DB.Preload("Photos").First(ps).RowsAffected < 1 {
		return "photoset does not exist", errors.New("not found")
	}
	org := &model.Organization{}
	if err := ctx.DB.First(org, "id = ?", orgID).Error; err != nil {
		return "target organization not found", err
	}
	store, err := ctx.StorageResolver.Storage(photosetLibraryID(ctx, ps))
	if err != nil {
		return "unable to resolve storage", err
	}
	for i := range ps.Photos {
		photo := &ps.Photos[i]
		srcPath := photo.RelativePath(ps)
		dstPath := org.RenderPhotoset(ps, photo)
		if srcPath == dstPath {
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
		_ = store.Delete(srcPath)
		path := dstPath
		photo.Path = &path
		if err := ctx.DB.Save(photo).Error; err != nil {
			return "unable to update photo path", err
		}
	}
	return "", nil
}

func reorganizeUpdateDB(ctx *common.Context, params common.Parameters) (string, error) {
	psID := params["photosetID"]
	orgID := params["targetOrganizationID"]
	ps := &model.Photoset{ID: psID}
	if ctx.DB.First(ps).RowsAffected < 1 {
		return "photoset does not exist", errors.New("not found")
	}
	if err := ctx.DB.Model(ps).Update("organization_id", orgID).Error; err != nil {
		return "unable to update photoset", err
	}
	return "", nil
}
