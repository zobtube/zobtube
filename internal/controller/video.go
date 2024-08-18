package controller

import (
	"errors"
	"fmt"
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

	"gitlab.com/zobtube/zobtube/internal/model"
)

var user = &model.User{
	Username: "test-user-admin",
	Admin:    true,
}

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
	c.datastore.Save(&video)

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

	// set new size
	dst := image.NewRGBA(image.Rect(0, 0, 320, 180))

	// resize
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

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
		"User":   user,
		"Video":  video,
	})
}

func (c *Controller) ClipList(g *gin.Context) {
	c.GenericVideoList("clip", g)
}

func (c *Controller) MovieList(g *gin.Context) {
	c.GenericVideoList("movie", g)
}

func (c *Controller) VideoList(g *gin.Context) {
	c.GenericVideoList("video", g)
}

func (c *Controller) GenericVideoList(videoType string, g *gin.Context) {
	var videos []model.Video
	c.datastore.Where("type = ?", videoType[0:1]).Order("created_at desc").Find(&videos)
	g.HTML(http.StatusOK, "video/list.html", gin.H{
		"Type":   videoType,
		"User":   user,
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
	c.datastore.Limit(12).Find(&randomVideos) //TODO: order by rand

	g.HTML(http.StatusOK, "video/view.html", gin.H{
		"Type":         video.Type,
		"User":         user,
		"Video":        video,
		"RandomVideos": randomVideos,
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
	fmt.Println(previousPath, " -> ", newPath)
	err = os.Rename(previousPath, newPath)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// commit the update on database
	video.Imported = true
	c.datastore.Save(video)
	//TODO: check result
	g.JSON(200, gin.H{})
}

func (c *Controller) VideoAjaxCreate(g *gin.Context) {
	var err error

	form := struct {
		ID           string   `form:"id"`
		Name         string   `form:"name"`
		Filename     string   `form:"filename"`
		Actors       []string `form:"actors"`
		TypeEnum     string   `form:"type"`
		CreationDate string   `form:"date"`
	}{}
	err = g.ShouldBind(&form)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// validate date
	date, err := time.Parse(time.RFC3339, form.CreationDate)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	video := &model.Video{
		ID:        form.ID,
		Name:      form.Name,
		Filename:  form.Filename,
		Type:      form.TypeEnum,
		CreatedAt: date,
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

	g.JSON(200, gin.H{})
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
