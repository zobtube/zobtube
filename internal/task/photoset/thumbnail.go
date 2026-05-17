package photoset

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
	"github.com/zobtube/zobtube/internal/task/common"
	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

const photoThumbMaxSide = 400

func generatePhotoMiniThumbnail(ctx *common.Context, ps *model.Photoset, photo *model.Photo) (string, error) {
	libStore, err := ctx.StorageResolver.Storage(photosetLibraryID(ctx, ps))
	if err != nil {
		return "unable to resolve library storage", err
	}
	rc, err := libStore.Open(photo.RelativePath(ps))
	if err != nil {
		return "unable to open photo", err
	}
	defer rc.Close()
	src, err := decodeImage(rc, photo.Filename)
	if err != nil {
		return "unable to decode image", err
	}
	dst := resizeToMaxSide(src, photoThumbMaxSide)
	thumbPath := photo.ThumbnailMiniRelativePath(ps)
	writeStore, err := metadataStoreForWrite(ctx)
	if err != nil {
		return "unable to resolve metadata storage", err
	}
	if err := writeStore.MkdirAll(filepath.Dir(thumbPath)); err != nil {
		return "unable to create thumbnail folder", err
	}
	w, err := writeStore.Create(thumbPath)
	if err != nil {
		return "unable to create mini thumbnail file", err
	}
	defer w.Close()
	if err := jpeg.Encode(w, dst, &jpeg.Options{Quality: 85}); err != nil {
		return "unable to encode mini thumbnail", err
	}
	photo.ThumbnailMini = true
	bounds := src.Bounds()
	photo.Width = bounds.Dx()
	photo.Height = bounds.Dy()
	if err := ctx.DB.Save(photo).Error; err != nil {
		return "unable to update photo row", err
	}
	return "", nil
}

func decodeImage(r io.Reader, filename string) (image.Image, error) {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Decode(r)
	case ".png":
		return png.Decode(r)
	case ".gif":
		return gif.Decode(r)
	case ".webp":
		return webp.Decode(r)
	default:
		img, _, err := image.Decode(r)
		return img, err
	}
}

func resizeToMaxSide(src image.Image, maxSide int) *image.RGBA {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= 0 || h <= 0 {
		return image.NewRGBA(image.Rect(0, 0, 1, 1))
	}
	scale := float64(maxSide) / float64(max(w, h))
	if scale >= 1 {
		dst := image.NewRGBA(image.Rect(0, 0, w, h))
		draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
		return dst
	}
	nw := int(float64(w) * scale)
	nh := int(float64(h) * scale)
	if nw < 1 {
		nw = 1
	}
	if nh < 1 {
		nh = 1
	}
	dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, b, draw.Over, nil)
	return dst
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func metadataStoreForWrite(ctx *common.Context) (storage.Storage, error) {
	if ctx.MetadataStorage != nil {
		return ctx.MetadataStorage, nil
	}
	return nil, errors.New("metadata storage not configured")
}
