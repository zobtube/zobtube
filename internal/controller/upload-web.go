package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) UploadTriage(g *gin.Context) {
	// get actors
	var actors []model.Actor
	err := c.datastore.Find(&actors).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// get categories
	categories := []model.Category{}
	err = c.datastore.Preload("Sub").Find(&categories).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// get channels
	var channels []model.Channel
	err = c.datastore.Find(&channels).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(g, http.StatusOK, "upload/home.html", gin.H{
		"Actors":     actors,
		"Categories": categories,
		"Channels":   channels,
	})
}

type UploadImportForm struct {
	Path     string `form:"path"`
	ImportAs string `form:"import_as"`
}

func (c *Controller) UploadImport(g *gin.Context) {
	var form UploadImportForm
	err := g.ShouldBind(&form)
	if err != nil {
		g.Redirect(http.StatusBadRequest, "/upload/triage")
		return
	}

	video := &model.Video{
		Name:          form.Path,
		Filename:      form.Path,
		Thumbnail:     false,
		ThumbnailMini: false,
		Type:          form.ImportAs,
	}

	c.datastore.Create(video)
	//TODO: check result
	g.Redirect(http.StatusFound, "/video/"+video.ID)
}
