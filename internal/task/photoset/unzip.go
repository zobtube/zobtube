package photoset

import (
	"archive/zip"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
	"github.com/zobtube/zobtube/internal/task/common"
)

func NewPhotosetUnzip() *common.Task {
	finalize := NewPhotosetFinalize()
	steps := []common.Step{
		{Name: "extract-archive", NiceName: "Extract archive", Func: extractArchive},
	}
	steps = append(steps, finalize.Steps...)
	return &common.Task{
		Name:  "photoset/unzip",
		Steps: steps,
	}
}

func extractArchive(ctx *common.Context, params common.Parameters) (string, error) {
	psID := params["photosetID"]
	archivePath := params["archivePath"]
	if archivePath == "" {
		return "archivePath required", errors.New("missing archivePath")
	}
	ps := &model.Photoset{ID: psID}
	if ctx.DB.First(ps).RowsAffected < 1 {
		return "photoset does not exist", errors.New("not found")
	}
	store, err := ctx.StorageResolver.Storage(photosetLibraryID(ctx, ps))
	if err != nil {
		return "unable to resolve storage", err
	}
	localPath, cleanup, err := storage.LocalPathForRead(store, archivePath)
	if err != nil {
		return "unable to read archive", err
	}
	defer cleanup()

	r, err := zip.OpenReader(localPath)
	if err != nil {
		return "unable to open zip", err
	}
	defer r.Close()

	imported := 0
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if strings.Contains(strings.ToLower(f.Name), "__macosx") {
			continue
		}
		base := sanitizeBasename(f.Name)
		if base == "" || strings.HasPrefix(base, ".") {
			continue
		}
		mimeType, ok := photoMimeForFilename(base)
		if !ok {
			continue
		}
		rel := filepath.Join(ps.FolderRelativePath(), base)
		rc, err := f.Open()
		if err != nil {
			return "unable to open zip entry", err
		}
		wc, err := store.Create(rel)
		if err != nil {
			rc.Close()
			return "unable to create file", err
		}
		_, err = io.Copy(wc, rc)
		rc.Close()
		_ = wc.Close()
		if err != nil {
			return "unable to write file", err
		}
		path := rel
		var existing model.Photo
		if ctx.DB.Where("photoset_id = ? AND filename = ?", ps.ID, base).First(&existing).RowsAffected > 0 {
			existing.Path = &path
			existing.Mime = mimeType
			existing.SizeBytes = int64(f.UncompressedSize64)
			_ = ctx.DB.Save(&existing).Error
			imported++
			continue
		}
		photo := &model.Photo{
			PhotosetID: ps.ID,
			Filename:   base,
			Path:       &path,
			Mime:       mimeType,
			SizeBytes:  int64(f.UncompressedSize64),
		}
		if err := ctx.DB.Create(photo).Error; err != nil {
			return "unable to create photo row", err
		}
		imported++
	}
	_ = store.Delete(archivePath)
	ps.Imported = true
	if err := ctx.DB.Save(ps).Error; err != nil {
		return "unable to update photoset", err
	}
	if imported == 0 {
		return "no supported images found in archive", errors.New("empty archive")
	}
	return "", nil
}
