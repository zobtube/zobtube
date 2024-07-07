package controller

import (
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"gitlab.com/zobtube/zobtube/internal/model"
)

func (c *Controller) UploadHome(g *gin.Context) {
	g.HTML(http.StatusOK, "upload/home.html", gin.H{
		"User": user,
	})
}

func (c *Controller) UploadTriage(g *gin.Context) {
	// list folders in triage path
	folders, err := os.ReadDir(
		filepath.Join(c.config.MediaFolder, TRIAGE_FILEPATH),
	)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	items := make(map[string][]fs.FileInfo)

	for _, folder := range folders {
		dir, err := os.Open(filepath.Join(
			c.config.MediaFolder,
			TRIAGE_FILEPATH,
			folder.Name(),
		))
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
		}
		defer dir.Close()

		// list files
		files, err := dir.Readdir(-1)
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
		}

		items[folder.Name()] = files
	}

	g.HTML(http.StatusOK, "upload/triage.html", gin.H{
		"User":  user,
		"Items": items,
	})
}

func (c *Controller) UploadPreview(g *gin.Context) {
	filePathEncoded := g.Param("filepath")
	filePath, err := url.QueryUnescape(filePathEncoded)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// construct file path
	targetPath := filepath.Join(c.config.MediaFolder, TRIAGE_FILEPATH, filePath)

	// give file path
	g.File(targetPath)
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
