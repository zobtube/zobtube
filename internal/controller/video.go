package controller

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/zobtube/zobtube/internal/model"
)

// VideoList godoc
//
//	@Summary	List all videos
//	@Tags		video
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Router		/video [get]
func (c *Controller) VideoList(g *gin.Context) {
	var videos []model.Video
	c.datastore.Where("type = ?", "v").Order("created_at desc").Find(&videos)
	g.JSON(http.StatusOK, gin.H{"items": videos, "total": len(videos)})
}

// ClipList godoc
//
//	@Summary	List all clips
//	@Tags		video
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Router		/clip [get]
func (c *Controller) ClipList(g *gin.Context) {
	var videos []model.Video
	c.datastore.Where("type = ?", "c").Order("created_at desc").Preload(clause.Associations).Find(&videos)
	g.JSON(http.StatusOK, gin.H{"items": videos, "total": len(videos)})
}

// MovieList godoc
//
//	@Summary	List all movies
//	@Tags		video
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Router		/movie [get]
func (c *Controller) MovieList(g *gin.Context) {
	var videos []model.Video
	c.datastore.Where("type = ?", "m").Order("created_at desc").Preload(clause.Associations).Find(&videos)
	g.JSON(http.StatusOK, gin.H{"items": videos, "total": len(videos)})
}

// VideoView godoc
//
//	@Summary	Get video view page data
//	@Tags		video
//	@Produce	json
//	@Param		id	path	string	true	"Video ID"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/video/{id} [get]
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

// VideoEdit godoc
//
//	@Summary	Get video edit form data (admin)
//	@Tags		video
//	@Produce	json
//	@Param		id	path	string	true	"Video ID"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/video/{id}/edit [get]
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

// VideoActors godoc
//
//	@Summary	Add or remove actor from video (PUT=add, DELETE=remove)
//	@Tags		video
//	@Param		id	path	string	true	"Video ID"
//	@Param		actor_id	path	string	true	"Actor ID"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/video/{id}/actor/{actor_id} [put]
//	@Router		/video/{id}/actor/{actor_id} [delete]
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

// VideoCategories godoc
//
//	@Summary	Add or remove category from video (PUT=add, DELETE=remove)
//	@Tags		video
//	@Param		id	path	string	true	"Video ID"
//	@Param		category_id	path	string	true	"Category (sub) ID"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/video/{id}/category/{category_id} [put]
//	@Router		/video/{id}/category/{category_id} [delete]
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

// VideoStreamInfo godoc
//
//	@Summary	Get video stream info (HEAD request)
//	@Tags		video
//	@Param		id	path	string	true	"Video ID"
//	@Success	200
//	@Failure	404
//	@Router		/video/{id} [head]
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

	libID := c.videoLibraryID(video)
	store, err := c.storageResolver.Storage(libID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	path := video.RelativePath()
	exists, err := store.Exists(path)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if exists {
		g.JSON(200, gin.H{})
	} else {
		g.JSON(404, gin.H{})
	}
}

type VideoRenameForm struct {
	Name string `form:"name"`
}

// VideoRename godoc
//
//	@Summary	Rename video
//	@Tags		video
//	@Accept		x-www-form-urlencoded
//	@Param		id	path	string	true	"Video ID"
//	@Param		name	formData	string	true	"New name"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/video/{id}/rename [post]
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

// VideoCreate godoc
//
//	@Summary	Create a new video
//	@Tags		video
//	@Accept		multipart/form-data
//	@Param		name	formData	string	true	"Video name"
//	@Param		filename	formData	string	true	"Filename in triage"
//	@Param		type	formData	string	true	"Type: c (clip), m (movie), v (video)"
//	@Param		actors	formData	array	false	"Actor IDs"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/video [post]
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

// VideoUploadThumb godoc
//
//	@Summary	Upload video thumbnail
//	@Tags		video
//	@Accept		multipart/form-data
//	@Param		id	path	string	true	"Video ID"
//	@Param		thumbnail	formData	file	true	"Thumbnail image"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/video/{id}/thumb [post]
func (c *Controller) VideoUploadThumb(g *gin.Context) {
	id := g.Param("id")
	video := &model.Video{ID: id}
	result := c.datastore.First(video)
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}
	store, err := c.storageResolver.Storage(c.videoLibraryID(video))
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := store.MkdirAll(video.FolderRelativePath()); err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	thumbnail, err := g.FormFile("thumbnail")
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	src, err := thumbnail.Open()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()
	dst, err := store.Create(video.ThumbnailRelativePath())
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	video.Thumbnail = true
	if err := c.datastore.Save(video).Error; err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := c.runner.NewTask("video/mini-thumb", map[string]string{"videoID": video.ID}); err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	g.JSON(200, gin.H{})
}

// VideoUpload godoc
//
//	@Summary	Upload video file
//	@Tags		video
//	@Accept		multipart/form-data
//	@Param		id	path	string	true	"Video ID"
//	@Param		file	formData	file	true	"Video file"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/video/{id}/upload [post]
func (c *Controller) VideoUpload(g *gin.Context) {
	id := g.Param("id")
	video := &model.Video{ID: id}
	result := c.datastore.First(video)
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}
	store, err := c.storageResolver.Storage(c.videoLibraryID(video))
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := store.MkdirAll(video.FolderRelativePath()); err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	videoData, err := g.FormFile("file")
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	src, err := videoData.Open()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()
	dst, err := store.Create(video.RelativePath())
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	video.Imported = true
	if err := c.datastore.Save(video).Error; err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	g.JSON(200, gin.H{})
}

// VideoDelete godoc
//
//	@Summary	Delete a video
//	@Tags		video
//	@Param		id	path	string	true	"Video ID"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/video/{id} [delete]
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

// VideoMigrate godoc
//
//	@Summary	Migrate video to different type (c/m/v)
//	@Tags		video
//	@Accept		x-www-form-urlencoded
//	@Param		id	path	string	true	"Video ID"
//	@Param		new_type	formData	string	true	"New type: c, m, or v"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/video/{id}/migrate [post]
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
	oldFolder := video.FolderRelativePathForType(video.Type)
	video.Type = newType
	newFolder := video.FolderRelativePath()

	store, err := c.storageResolver.Storage(c.videoLibraryID(video))
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Copy known files from old folder to new folder
	for _, name := range []string{"video.mp4", "thumb.jpg", "thumb-xs.jpg"} {
		oldPath := filepath.Join(oldFolder, name)
		newPath := filepath.Join(newFolder, name)
		exists, _ := store.Exists(oldPath)
		if !exists {
			continue
		}
		rc, err := store.Open(oldPath)
		if err != nil {
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}
		wc, err := store.Create(newPath)
		if err != nil {
			rc.Close()
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}
		_, err = io.Copy(wc, rc)
		rc.Close()
		wc.Close()
		if err != nil {
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}
		_ = store.Delete(oldPath)
	}

	if err := c.datastore.Save(video).Error; err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	g.JSON(200, gin.H{})
}

// VideoGenerateThumbnail godoc
//
//	@Summary	Generate video thumbnail at given timing
//	@Tags		video
//	@Param		id	path	string	true	"Video ID"
//	@Param		timing	path	string	true	"Timing in seconds"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Failure	409	{object}	map[string]interface{}
//	@Router		/video/{id}/generate-thumbnail/{timing} [post]
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

// VideoGet godoc
//
//	@Summary	Get video summary (title, actors, categories)
//	@Tags		video
//	@Produce	json
//	@Param		id	path	string	true	"Video ID"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	404
//	@Router		/video/{id}/summary [get]
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

// VideoStream godoc
//
//	@Summary	Stream video file
//	@Tags		video
//	@Param		id	path	string	true	"Video ID"
//	@Success	200	file	bytes
//	@Failure	404
//	@Router		/video/{id}/stream [get]
func (c *Controller) VideoStream(g *gin.Context) {
	id := g.Param("id")
	video := &model.Video{ID: id}
	result := c.datastore.First(video)
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}
	libID := c.videoLibraryID(video)
	store, err := c.storageResolver.Storage(libID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var path string
	if video.Imported {
		path = video.RelativePath()
	} else {
		path = filepath.Join("triage", video.Filename)
	}
	c.serveFromStorage(g, store, path)
}

// VideoThumb godoc
//
//	@Summary	Get video thumbnail image
//	@Tags		video
//	@Param		id	path	string	true	"Video ID"
//	@Success	200	file	bytes
//	@Failure	404
//	@Router		/video/{id}/thumb [get]
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
	store, err := c.storageResolver.Storage(c.videoLibraryID(video))
	if err != nil {
		c.logger.Error().Err(err).Str("video_id", video.ID).Str("library_id", c.videoLibraryID(video)).Msg("error resolving storage")
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("serving thumbnail from", video.ThumbnailRelativePath())
	c.serveFromStorage(g, store, video.ThumbnailRelativePath())
}

// VideoThumbXS godoc
//
//	@Summary	Get video extra-small thumbnail
//	@Tags		video
//	@Param		id	path	string	true	"Video ID"
//	@Success	200	file	bytes
//	@Failure	404
//	@Router		/video/{id}/thumb_xs [get]
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
	store, err := c.storageResolver.Storage(c.videoLibraryID(video))
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.serveFromStorage(g, store, video.ThumbnailXSRelativePath())
}
