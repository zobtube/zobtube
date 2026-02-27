package video

import (
	"errors"
	"image"
	"image/jpeg"
	"path/filepath"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
	"github.com/zobtube/zobtube/internal/task/common"
	"golang.org/x/image/draw"
)

func generateHorizontalMiniThumnail(ctx *common.Context, video *model.Video, store storage.Storage) (string, error) {
	thumbPath := video.ThumbnailRelativePath()
	thumbXSPath := video.ThumbnailXSRelativePath()
	rc, err := store.Open(thumbPath)
	if err != nil {
		return "unable to open thumbnail", err
	}
	defer rc.Close()
	src, err := jpeg.Decode(rc)
	if err != nil {
		return "unable to read the jpg file", err
	}
	targetH, targetV := 320, 180
	h := src.Bounds().Dx()
	v := src.Bounds().Dy()
	originalImageRGBA := image.NewRGBA(image.Rect(0, 0, h, v))
	draw.Draw(originalImageRGBA, originalImageRGBA.Bounds(), src, src.Bounds().Min, draw.Src)
	ratioH := float32(h) / float32(targetH)
	ratioV := float32(v) / float32(targetV)
	ratio := max(ratioH, ratioV)
	h = int(float32(h) / ratio)
	v = int(float32(v) / ratio)
	dst := image.NewRGBA(image.Rect(0, 0, targetH, targetV))
	outerImg := gaussianBlur(originalImageRGBA, 15)
	draw.NearestNeighbor.Scale(dst, dst.Bounds(), outerImg, outerImg.Bounds(), draw.Over, nil)
	innerH := (targetH - h) / 2
	innerV := (targetV - v) / 2
	draw.NearestNeighbor.Scale(dst, image.Rect(innerH, innerV, innerH+h, innerV+v), src, src.Bounds(), draw.Over, nil)
	if err := store.MkdirAll(filepath.Dir(thumbXSPath)); err != nil {
		return "unable to create thumbnail folder", err
	}
	w, err := store.Create(thumbXSPath)
	if err != nil {
		return "unable to create mini thumbnail file", err
	}
	defer w.Close()
	if err := jpeg.Encode(w, dst, &jpeg.Options{Quality: 90}); err != nil {
		return "unable to encode new thumbnail", err
	}
	return "", nil
}

func generateSameRatioMiniThumnail(ctx *common.Context, video *model.Video, store storage.Storage) (string, error) {
	thumbPath := video.ThumbnailRelativePath()
	thumbXSPath := video.ThumbnailXSRelativePath()
	rc, err := store.Open(thumbPath)
	if err != nil {
		return "unable to open thumbnail", err
	}
	defer rc.Close()
	src, err := jpeg.Decode(rc)
	if err != nil {
		return "unable to read the jpg file", err
	}
	targetH := 320
	h := src.Bounds().Dx()
	v := src.Bounds().Dy()
	var dst *image.RGBA
	if h <= targetH {
		dst = image.NewRGBA(image.Rect(0, 0, h, v))
		draw.NearestNeighbor.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	} else {
		ratio := float32(h) / float32(targetH)
		v = int(float32(v) / ratio)
		dst = image.NewRGBA(image.Rect(0, 0, targetH, v))
		draw.NearestNeighbor.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	}
	if err := store.MkdirAll(filepath.Dir(thumbXSPath)); err != nil {
		return "unable to create thumbnail folder", err
	}
	w, err := store.Create(thumbXSPath)
	if err != nil {
		return "unable to create mini thumbnail file", err
	}
	defer w.Close()
	if err := jpeg.Encode(w, dst, &jpeg.Options{Quality: 90}); err != nil {
		return "unable to encode new thumbnail", err
	}
	return "", nil
}

func generateThumbnailMini(ctx *common.Context, params common.Parameters) (string, error) {
	videoID := params["videoID"]
	video := &model.Video{ID: videoID}
	result := ctx.DB.First(video)
	if result.RowsAffected < 1 {
		return "video does not exist", errors.New("id not in db")
	}
	store, err := ctx.StorageResolver.Storage(videoLibraryID(ctx, video))
	if err != nil {
		return "unable to resolve storage", err
	}
	var errMsg string
	if video.Type == "c" {
		errMsg, err = generateSameRatioMiniThumnail(ctx, video, store)
	} else {
		errMsg, err = generateHorizontalMiniThumnail(ctx, video, store)
	}
	if err != nil {
		return errMsg, err
	}
	video.ThumbnailMini = true
	ctx.DB.Save(&video)
	return "", nil
}

func deleteThumbnailMini(ctx *common.Context, params common.Parameters) (string, error) {
	videoID := params["videoID"]
	video := &model.Video{ID: videoID}
	result := ctx.DB.First(video)
	if result.RowsAffected < 1 {
		return "video does not exist", errors.New("id not in db")
	}
	store, err := ctx.StorageResolver.Storage(videoLibraryID(ctx, video))
	if err != nil {
		return "unable to resolve storage", err
	}
	_ = store.Delete(video.ThumbnailXSRelativePath())
	return "", nil
}
