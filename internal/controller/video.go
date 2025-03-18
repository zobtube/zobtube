package controller

import (
	"errors"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/draw"
	"gorm.io/gorm/clause"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) VideoAjaxActors(g *gin.Context) {
	// get id from path
	id := g.Param("id")
	actor_id := g.Param("actor_id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	actor := &model.Actor{
		ID: actor_id,
	}
	result = c.datastore.First(&actor)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	var res error
	if g.Request.Method == "PUT" {
		res = c.datastore.Model(video).Association("Actors").Append(actor)
	} else {
		res = c.datastore.Model(video).Association("Actors").Delete(actor)
	}

	if res != nil {
		g.JSON(500, gin.H{})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) VideoAjaxComputeDuration(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	filePath := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], id, "video.mp4")
	out, err := exec.Command(
		"ffprobe",
		"-v",
		"error",
		"-show_entries",
		"format=duration",
		"-of",
		"default=noprint_wrappers=1:nokey=1",
		filePath,
	).Output()

	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
			"out":   string(out[:]),
		})
		return
	}

	duration := strings.TrimSpace(string(out))

	d, err := time.ParseDuration(duration + "s")
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	video.Duration = d
	err = c.datastore.Save(&video).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{
		"duration": d.String(),
	})
}

func (c *Controller) VideoAjaxGenerateThumbnail(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get timing from path
	timing := g.Param("timing")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// construct paths
	videoPath := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], id, "video.mp4")
	videoPath, err := filepath.Abs(videoPath)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	thumbPath := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], id, "thumb.jpg")
	thumbPath, err = filepath.Abs(thumbPath)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	out, err := exec.Command(
		"ffmpeg",
		"-y",
		"-ss",
		timing,
		"-i",
		videoPath,
		"-frames:v",
		"1",
		"-q:v",
		"2",
		thumbPath,
	).Output()
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
			"out":   string(out[:]),
		})
		return
	}

	video.Thumbnail = true
	c.datastore.Save(&video)

	g.JSON(200, gin.H{})
}

func (c *Controller) VideoAjaxGenerateThumbnailXS(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// construct paths
	thumbPath := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], id, "thumb.jpg")
	thumbPath, err := filepath.Abs(thumbPath)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	thumbXSPath := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], id, "thumb-xs.jpg")
	thumbXSPath, err = filepath.Abs(thumbXSPath)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// open files
	input, _ := os.Open(thumbPath)
	defer input.Close()

	output, _ := os.Create(thumbXSPath)
	defer output.Close()

	// decode the image from jpeg to image.Image
	src, err := jpeg.Decode(input)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
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
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// save on db
	video.ThumbnailMini = true
	c.datastore.Save(&video)

	// ret
	g.JSON(200, gin.H{})
}

func (c *Controller) VideoEdit(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.Preload("Actors").First(video)

	var actors []model.Actor
	c.datastore.Find(&actors)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	g.HTML(http.StatusOK, "video/edit.html", gin.H{
		"Actors": actors,
		"User":   g.MustGet("user").(*model.User),
		"Video":  video,
	})
}

func (c *Controller) ClipList(g *gin.Context) {
	c.GenericVideoList("clip", g)
}

func (c *Controller) MovieList(g *gin.Context) {
	var videos []model.Video
	c.datastore.Where("type = ?", "m").Order("created_at desc").Preload(clause.Associations).Find(&videos)
	g.HTML(http.StatusOK, "movie/list.html", gin.H{
		"Type":   "movie",
		"User":   g.MustGet("user").(*model.User),
		"Videos": videos,
	})
}

func (c *Controller) VideoList(g *gin.Context) {
	c.GenericVideoList("video", g)
}

func (c *Controller) GenericVideoList(videoType string, g *gin.Context) {
	var videos []model.Video
	c.datastore.Where("type = ?", videoType[0:1]).Order("created_at desc").Find(&videos)
	g.HTML(http.StatusOK, "video/list.html", gin.H{
		"Type":   videoType,
		"User":   g.MustGet("user").(*model.User),
		"Videos": videos,
	})
}

func (c *Controller) VideoView(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.Preload("Actors").First(video)

	// check result
	if result.RowsAffected < 1 {
		//TODO: return to homepage
		g.JSON(404, gin.H{})
		return
	}

	// get random videos
	var randomVideos []model.Video
	c.datastore.Limit(8).Order("RANDOM()").Find(&randomVideos)

	// get video count
	user := g.MustGet("user").(*model.User)
	viewCount := 0
	count := &model.VideoView{}
	result = c.datastore.First(&count, "video_id = ? AND user_id = ?", video.ID, user.ID)
	if result.RowsAffected > 0 {
		viewCount = count.Count
	}

	g.HTML(http.StatusOK, "video/view.html", gin.H{
		"Type":      video.Type,
		"User":      user,
		"Video":     video,
		"ViewCount": viewCount,
		"RandomVideos": gin.H{
			"Videos":    randomVideos,
			"VideoType": video.Type,
		},
	})
}

func (c *Controller) VideoStream(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// construct file path
	var targetPath string
	if video.Imported {
		targetPath = filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], id, "video.mp4")
	} else {
		targetPath = filepath.Join(c.config.Media.Path, TRIAGE_FILEPATH, video.Filename)
	}

	// give file path
	g.File(targetPath)
}

func (c *Controller) VideoAjaxStreamInfo(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	if !video.Imported {
		g.JSON(404, gin.H{})
		return
	}

	path := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], video.ID, "video.mp4")
	_, err := os.Stat(path)
	if err == nil {
		g.JSON(200, gin.H{})
		return
	} else if errors.Is(err, os.ErrNotExist) {
		g.JSON(404, gin.H{})
		return
	} else {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func (c *Controller) VideoThumb(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// check if thumbnail exists
	if !video.Thumbnail {
		g.JSON(404, gin.H{})
		return
	}

	// construct file path
	targetPath := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], id, "thumb.jpg")

	// give file path
	g.File(targetPath)
}

func (c *Controller) VideoThumbXS(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// check if thumbnail exists
	if !video.ThumbnailMini {
		g.Redirect(http.StatusFound, VIDEO_THUMB_NOT_GENERATED)
		return
	}

	// construct file path
	targetPath := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], id, "thumb-xs.jpg")

	// give file path
	g.File(targetPath)
}

type VideoAjaxRenameForm struct {
	Name string `form:"name"`
}

func (c *Controller) VideoAjaxRename(g *gin.Context) {
	if g.Request.Method != "POST" {
		// method not allowed
		g.JSON(405, gin.H{})
		return
	}

	var form VideoAjaxRenameForm
	err := g.ShouldBind(&form)
	if err != nil {
		// method not allowed
		g.JSON(406, gin.H{})
		return
	}

	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	video.Name = form.Name

	c.datastore.Save(video)
	//TODO: check result
	g.JSON(200, gin.H{})
}

func (c *Controller) VideoAjaxImport(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// prepare paths
	previousPath := filepath.Join(c.config.Media.Path, TRIAGE_FILEPATH, video.Filename)
	newFolderPath := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], id)
	newPath := filepath.Join(newFolderPath, "video.mp4")

	// ensure folder exists
	_, err := os.Stat(newFolderPath)
	if os.IsNotExist(err) {
		// do not exists, create it
		err = os.Mkdir(newFolderPath, os.ModePerm)
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// move
	err = os.Rename(previousPath, newPath)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// commit the update on database
	video.Imported = true
	err = c.datastore.Save(video).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) VideoAjaxCreate(g *gin.Context) {
	var err error

	form := struct {
		Name     string   `form:"name"`
		Filename string   `form:"filename"`
		Actors   []string `form:"actors"`
		TypeEnum string   `form:"type"`
	}{}
	err = g.ShouldBind(&form)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	if form.Name == "" || form.Filename == "" || (form.TypeEnum != "c" && form.TypeEnum != "m" && form.TypeEnum != "v") {
		g.JSON(500, gin.H{
			"error": "invalid input",
		})
		return
	}

	video := &model.Video{
		Name:      form.Name,
		Filename:  form.Filename,
		Type:      form.TypeEnum,
		Imported:  false,
		Thumbnail: false,
	}

	// save object in db
	err = c.datastore.Create(&video).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, actorID := range form.Actors {
		actor := &model.Actor{
			ID: actorID,
		}
		result := c.datastore.First(&actor)

		// check result
		if result.RowsAffected < 1 {
			g.JSON(404, gin.H{})
			return
		}

		err = c.datastore.Model(video).Association("Actors").Append(actor)
		if err != nil {
			c.datastore.Delete(&video)
			g.JSON(500, gin.H{})
			return
		}
	}

	g.JSON(200, gin.H{
		"video_id": video.ID,
	})
}

func (c *Controller) VideoAjaxUploadThumb(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// ensure folder exists
	videoFolder := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], video.ID)
	_, err := os.Stat(videoFolder)
	if os.IsNotExist(err) {
		// do not exists, create it
		err = os.Mkdir(videoFolder, os.ModePerm)
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// save thumbnail
	thumbnailPath := filepath.Join(videoFolder, "thumb.jpg")
	thumbnail, err := g.FormFile("thumbnail")
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = g.SaveUploadedFile(thumbnail, thumbnailPath)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// commit the update on database
	video.Thumbnail = true
	err = c.datastore.Save(video).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// generate mini thumb
	c.VideoAjaxGenerateThumbnailXS(g)
}

func (c *Controller) VideoAjaxUpload(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// ensure folder exists
	videoFolder := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], video.ID)
	_, err := os.Stat(videoFolder)
	if os.IsNotExist(err) {
		// do not exists, create it
		err = os.Mkdir(videoFolder, os.ModePerm)
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// save video
	videoPath := filepath.Join(videoFolder, "video.mp4")
	videoData, err := g.FormFile("file")
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = g.SaveUploadedFile(videoData, videoPath)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// commit the update on database
	video.Imported = true
	err = c.datastore.Save(video).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) VideoAjaxDelete(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// check thumb presence
	thumbPath := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], video.ID, "thumb.jpg")
	_, err := os.Stat(thumbPath)
	if err != nil && !os.IsNotExist(err) {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !os.IsNotExist(err) {
		// exist, deleting it
		err = os.Remove(thumbPath)
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// check thumb-xs presence
	thumbXsPath := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], video.ID, "thumb-xs.jpg")
	_, err = os.Stat(thumbXsPath)
	if err != nil && !os.IsNotExist(err) {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !os.IsNotExist(err) {
		// exist, deleting it
		err = os.Remove(thumbXsPath)
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// check video presence
	videoPath := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], video.ID, "video.mp4")
	_, err = os.Stat(videoPath)
	if err != nil && !os.IsNotExist(err) {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !os.IsNotExist(err) {
		// exist, deleting it
		err = os.Remove(videoPath)
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// delete folder
	folderPath := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], video.ID)
	_, err = os.Stat(folderPath)
	if err != nil && !os.IsNotExist(err) {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !os.IsNotExist(err) {
		// exist, deleting it
		err = os.Remove(folderPath)
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// delete object
	err = c.datastore.Delete(video).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) VideoAjaxMigrate(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	newType := g.PostForm("new_type")
	newVid := &model.Video{
		Type: newType,
	}

	previousPath := filepath.Join(c.config.Media.Path, fileTypeToPath[video.TypeAsString()], video.ID)
	newPath := filepath.Join(c.config.Media.Path, fileTypeToPath[newVid.TypeAsString()], video.ID)

	// move
	err := os.Rename(previousPath, newPath)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// change object in db
	video.Type = newType

	// commit
	err = c.datastore.Save(video).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{})
}
