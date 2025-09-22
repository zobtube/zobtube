package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) UploadTriage(g *gin.Context) {
	c.HTML(g, http.StatusOK, "upload/home.html", gin.H{})
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
