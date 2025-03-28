package controller

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

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

	path := filepath.Join(c.config.Media.Path, video.RelativePath())
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
	videoFolder := filepath.Join(c.config.Media.Path, video.FolderRelativePath())
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
	videoFolder := filepath.Join(c.config.Media.Path, video.FolderRelativePath())
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

	previousPath := filepath.Join(c.config.Media.Path, video.RelativePath())

	// change object in db
	video.Type = newType

	newPath := filepath.Join(c.config.Media.Path, video.RelativePath())

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
