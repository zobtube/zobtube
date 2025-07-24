package controller

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) VideoEdit(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.Preload("Actors").Preload("Channel").Preload("Categories").First(video)

	var actors []model.Actor
	c.datastore.Find(&actors)

	// get categories
	categories := []model.Category{}
	err := c.datastore.Preload("Sub").Find(&categories).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	g.HTML(http.StatusOK, "video/edit.html", gin.H{
		"Actors":     actors,
		"User":       g.MustGet("user").(*model.User),
		"Video":      video,
		"Categories": categories,
	})
}

func (c *Controller) ClipList(g *gin.Context) {
	var videos []model.Video
	c.datastore.Where("type = ?", "c").Order("created_at desc").Preload(clause.Associations).Find(&videos)
	g.HTML(http.StatusOK, "clip/list.html", gin.H{
		"Type":   "clip",
		"User":   g.MustGet("user").(*model.User),
		"Videos": videos,
	})
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
	result := c.datastore.Preload("Actors.Categories").Preload("Channel").Preload("Categories").First(video)

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

	// get categories
	categories := make(map[string]string)
	for _, category := range video.Categories {
		categories[category.ID] = category.Name
	}
	for _, actor := range video.Actors {
		for _, category := range actor.Categories {
			categories[category.ID] = category.Name
		}
	}

	g.HTML(http.StatusOK, "video/view.html", gin.H{
		"Type":       video.Type,
		"User":       user,
		"Video":      video,
		"ViewCount":  viewCount,
		"Categories": categories,
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
		targetPath = filepath.Join(c.config.Media.Path, video.RelativePath())
	} else {
		targetPath = filepath.Join(c.config.Media.Path, TRIAGE_FILEPATH, video.Filename)
	}

	// give file path
	g.File(targetPath)
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
	targetPath := filepath.Join(c.config.Media.Path, video.ThumbnailRelativePath())

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
	targetPath := filepath.Join(c.config.Media.Path, video.ThumbnailXSRelativePath())

	// give file path
	g.File(targetPath)
}
