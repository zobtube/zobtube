package video

import (
	"errors"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/task/common"
	"golang.org/x/image/draw"
)

func generateThumbnailMini(ctx *common.Context, params common.Parameters) (string, error) {
	videoID := params["videoID"]

	// get item from ID
	video := &model.Video{
		ID: videoID,
	}
	result := ctx.DB.First(video)

	// check result
	if result.RowsAffected < 1 {
		return "video does not exist", errors.New("id not in db")
	}

	// construct paths
	thumbPath := filepath.Join(ctx.Config.Media.Path, video.ThumbnailRelativePath())
	thumbPath, err := filepath.Abs(thumbPath)
	if err != nil {
		return "Unable to get absolute path of the thumbnail", err
	}

	thumbXSPath := filepath.Join(ctx.Config.Media.Path, video.ThumbnailXSRelativePath())
	thumbXSPath, err = filepath.Abs(thumbXSPath)
	if err != nil {
		return "Unable to get absolute path of the new mini thumbnail", err
	}

	// open files
	input, _ := os.Open(thumbPath)
	defer input.Close()

	output, _ := os.Create(thumbXSPath)
	defer output.Close()

	// decode the image from jpeg to image.Image
	src, err := jpeg.Decode(input)
	if err != nil {
		return "unable to read the jpg file", err
	}

	targetH := 320
	targetV := 180

	h := src.Bounds().Dx()
	v := src.Bounds().Dy()

	originalImageRGBA := image.NewRGBA(image.Rect(0, 0, h, v))
	draw.Draw(originalImageRGBA, originalImageRGBA.Bounds(), src, src.Bounds().Min, draw.Src)

	ratioH := float32(h) / float32(targetH)
	ratioV := float32(v) / float32(targetV)
	ratio := max(ratioH, ratioV)

	h = int(float32(h) / ratio)
	v = int(float32(v) / ratio)

	// set new size
	dst := image.NewRGBA(image.Rect(0, 0, targetH, targetV))

	// draw outer
	outerImg := gaussianBlur(originalImageRGBA, 15)
	draw.NearestNeighbor.Scale(dst, dst.Bounds(), outerImg, outerImg.Bounds(), draw.Over, nil)

	// draw inner
	innerH := (targetH - h) / 2
	innerV := (targetV - v) / 2
	draw.NearestNeighbor.Scale(dst, image.Rect(innerH, innerV, innerH+h, innerV+v), src, src.Bounds(), draw.Over, nil)

	// encode to jpeg
	err = jpeg.Encode(output, dst, &jpeg.Options{Quality: 90})
	if err != nil {
		return "unable to encode new thumbnail", err
	}

	// save on db
	video.ThumbnailMini = true
	ctx.DB.Save(&video)

	// ret
	return "", nil
}

func deleteThumbnailMini(ctx *common.Context, params common.Parameters) (string, error) {
	videoID := params["videoID"]

	// get item from ID
	video := &model.Video{
		ID: videoID,
	}
	result := ctx.DB.First(video)

	// check result
	if result.RowsAffected < 1 {
		return "video does not exist", errors.New("id not in db")
	}

	// check thumb-xs presence
	thumbXsPath := filepath.Join(ctx.Config.Media.Path, video.ThumbnailXSRelativePath())
	_, err := os.Stat(thumbXsPath)
	if err != nil && !os.IsNotExist(err) {
		return "unable to check mini thumbnail presence", err
	}
	if !os.IsNotExist(err) {
		// exist, deleting it
		err = os.Remove(thumbXsPath)
		if err != nil {
			return "unable to delete mini thumbnail", err
		}
	}

	return "", nil
}
