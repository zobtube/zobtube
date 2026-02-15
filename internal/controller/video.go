package controller

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) VideoList(g *gin.Context) {
	var videos []model.Video
	c.datastore.Where("type = ?", "v").Order("created_at desc").Find(&videos)
	g.JSON(http.StatusOK, gin.H{"items": videos, "total": len(videos)})
}

func (c *Controller) ClipList(g *gin.Context) {
	var videos []model.Video
	c.datastore.Where("type = ?", "c").Order("created_at desc").Preload(clause.Associations).Find(&videos)
	g.JSON(http.StatusOK, gin.H{"items": videos, "total": len(videos)})
}

func (c *Controller) MovieList(g *gin.Context) {
	var videos []model.Video
	c.datastore.Where("type = ?", "m").Order("created_at desc").Preload(clause.Associations).Find(&videos)
	g.JSON(http.StatusOK, gin.H{"items": videos, "total": len(videos)})
}

func (c *Controller) VideoView(g *gin.Context) {
	id := g.Param("id")
	video := &model.Video{ID: id}
	result := c.datastore.Preload("Actors.Categories").Preload("Channel").Preload("Categories").First(video)
	if result.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	viewCount := 0
	user, ok := g.Get("user")
	if ok {
		if u, ok := user.(*model.User); ok && u != nil && u.ID != "" {
			count := &model.VideoView{}
			if c.datastore.First(count, "video_id = ? AND user_id = ?", video.ID, u.ID).RowsAffected > 0 {
				viewCount = count.Count
			}
		}
	}
	categories := make(map[string]string)
	for _, cat := range video.Categories {
		categories[cat.ID] = cat.Name
	}
	for _, actor := range video.Actors {
		for _, cat := range actor.Categories {
			categories[cat.ID] = cat.Name
		}
	}
	var randomVideos []model.Video
	c.datastore.Limit(8).Where("type = ? and id != ?", video.Type, video.ID).Order("RANDOM()").Find(&randomVideos)
	g.JSON(http.StatusOK, gin.H{
		"video":         video,
		"view_count":    viewCount,
		"categories":    categories,
		"random_videos": randomVideos,
	})
}

func (c *Controller) VideoEdit(g *gin.Context) {
	id := g.Param("id")
	video := &model.Video{ID: id}
	result := c.datastore.Preload("Actors").Preload("Channel").Preload("Categories").First(video)
	if result.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var actors []model.Actor
	c.datastore.Find(&actors)
	var categories []model.Category
	c.datastore.Preload("Sub").Find(&categories)
	g.JSON(http.StatusOK, gin.H{
		"video":      video,
		"actors":     actors,
		"categories": categories,
	})
}

func (c *Controller) VideoActors(g *gin.Context) {
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

func (c *Controller) VideoCategories(g *gin.Context) {
	// get id from path
	id := g.Param("id")
	category_id := g.Param("category_id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{
			"error": "video not found",
		})
		return
	}

	subCategory := &model.CategorySub{
		ID: category_id,
	}
	result = c.datastore.First(&subCategory)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{
			"error": "sub-category not found",
		})
		return
	}

	var res error
	if g.Request.Method == "PUT" {
		res = c.datastore.Model(video).Association("Categories").Append(subCategory)
	} else {
		res = c.datastore.Model(video).Association("Categories").Delete(subCategory)
	}

	if res != nil {
		g.JSON(500, gin.H{
			"error": res.Error(),
		})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) VideoStreamInfo(g *gin.Context) {
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

	path := filepath.Join(c.config.Media.Path, video.RelativePath())
	_, err := os.Stat(path)
	if err == nil {
		g.JSON(200, gin.H{})
	} else if errors.Is(err, os.ErrNotExist) {
		g.JSON(404, gin.H{})
	} else {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
}

type VideoRenameForm struct {
	Name string `form:"name"`
}

func (c *Controller) VideoRename(g *gin.Context) {
	if g.Request.Method != "POST" {
		// method not allowed
		g.JSON(405, gin.H{})
		return
	}

	var form VideoRenameForm
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

	err = c.datastore.Save(video).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) VideoCreate(g *gin.Context) {
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

	err = c.runner.NewTask("video/create", map[string]string{
		"videoID":         video.ID,
		"thumbnailTiming": "0",
	})
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, gin.H{
		"video_id": video.ID,
	})
}

func (c *Controller) VideoUploadThumb(g *gin.Context) {
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
	videoFolder := filepath.Join(c.config.Media.Path, video.FolderRelativePath())
	_, err := os.Stat(videoFolder)
	if os.IsNotExist(err) {
		// do not exists, create it
		err = os.Mkdir(videoFolder, 0o750)
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
	thumbnailPath := video.ThumbnailRelativePath()
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

	err = c.runner.NewTask("video/mini-thumb", map[string]string{"videoID": video.ID})
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) VideoUpload(g *gin.Context) {
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
	videoFolder := filepath.Join(c.config.Media.Path, video.FolderRelativePath())
	_, err := os.Stat(videoFolder)
	if os.IsNotExist(err) {
		// do not exists, create it
		err = os.Mkdir(videoFolder, 0o750)
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

func (c *Controller) VideoDelete(g *gin.Context) {
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

	// update status
	video.Status = model.VideoStatusDeleting
	err := c.datastore.Save(video).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// create task
	err = c.runner.NewTask("video/delete", map[string]string{
		"videoID": video.ID,
	})
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) VideoMigrate(g *gin.Context) {
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

	previousPath := filepath.Join(c.config.Media.Path, video.FolderRelativePath())

	// change object in db
	video.Type = newType

	newPath := filepath.Join(c.config.Media.Path, video.FolderRelativePath())

	// move
	err := os.Rename(previousPath, newPath)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

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

func (c *Controller) VideoGenerateThumbnail(g *gin.Context) {
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

	if video.Status != model.VideoStatusReady {
		g.JSON(409, gin.H{
			"error": "video is not ready to be updated",
		})
		return
	}

	// create task
	err := c.runner.NewTask("video/generate-thumbnail", map[string]string{
		"videoID":         video.ID,
		"thumbnailTiming": g.Param("timing"),
	})
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) VideoEditChannel(g *gin.Context) {
	if g.Request.Method != "POST" {
		// method not allowed
		g.JSON(405, gin.H{})
		return
	}

	form := struct {
		ChannelID string `form:"channelID"`
	}{}

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

	if form.ChannelID == "x" {
		video.ChannelID = nil
	} else {
		channel := &model.Channel{
			ID: form.ChannelID,
		}

		result := c.datastore.First(channel)
		// check result
		if result.RowsAffected < 1 {
			g.JSON(404, gin.H{})
			return
		}

		video.Channel = channel
	}

	err = c.datastore.Save(video).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) VideoGet(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.Preload("Actors.Categories").Preload("Categories").First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	var actors []string
	var categories []string

	for _, actor := range video.Actors {
		actors = append(actors, actor.Name)
		for _, category := range actor.Categories {
			categories = append(categories, category.Name)
		}
	}

	for _, category := range video.Categories {
		categories = append(categories, category.Name)
	}

	g.JSON(200, gin.H{
		"title":      video.Name,
		"actors":     actors,
		"categories": categories,
	})
}

func (c *Controller) VideoStream(g *gin.Context) {
	id := g.Param("id")
	video := &model.Video{ID: id}
	result := c.datastore.First(video)
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}
	var targetPath string
	if video.Imported {
		targetPath = filepath.Join(c.config.Media.Path, video.RelativePath())
	} else {
		targetPath = filepath.Join(c.config.Media.Path, TRIAGE_FILEPATH, video.Filename)
	}
	g.File(targetPath)
}

func (c *Controller) VideoThumb(g *gin.Context) {
	id := g.Param("id")
	video := &model.Video{ID: id}
	result := c.datastore.First(video)
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}
	if !video.Thumbnail {
		g.JSON(404, gin.H{})
		return
	}
	targetPath := filepath.Join(c.config.Media.Path, video.ThumbnailRelativePath())
	g.File(targetPath)
}

func (c *Controller) VideoThumbXS(g *gin.Context) {
	id := g.Param("id")
	video := &model.Video{ID: id}
	result := c.datastore.First(video)
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}
	if !video.ThumbnailMini {
		g.Redirect(http.StatusFound, VIDEO_THUMB_NOT_GENERATED)
		return
	}
	targetPath := filepath.Join(c.config.Media.Path, video.ThumbnailXSRelativePath())
	g.File(targetPath)
}
