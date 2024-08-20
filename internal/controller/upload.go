package controller

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"gitlab.com/zobtube/zobtube/internal/model"
)

func (c *Controller) UploadHome(g *gin.Context) {
	g.HTML(http.StatusOK, "upload/home.html", gin.H{
		"User": g.MustGet("user").(*model.User),
	})
}

func (c *Controller) UploadTriage(g *gin.Context) {
	g.HTML(http.StatusOK, "upload/triage.html", gin.H{
		"User": g.MustGet("user").(*model.User),
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
	targetPath := filepath.Join(c.config.Media.Path, TRIAGE_FILEPATH, filePath)

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

func (c *Controller) UploadAjaxTriageFolder(g *gin.Context) {
	// get requested path
	path := g.PostForm("path")

	// list folders in triage path
	folders, err := os.ReadDir(
		filepath.Join(c.config.Media.Path, TRIAGE_FILEPATH, path),
	)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	items := make(map[string]int)

	for _, folder := range folders {
		// check type
		entryPath := filepath.Join(
			c.config.Media.Path,
			TRIAGE_FILEPATH,
			path,
			folder.Name(),
		)
		stat, err := os.Stat(entryPath)
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		if !stat.IsDir() {
			continue
		}

		dir, err := os.Open(entryPath)
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer dir.Close()

		// list files
		files, err := dir.Readdir(-1)
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		items[folder.Name()] = len(files)
	}

	g.JSON(http.StatusOK, gin.H{
		"folders": items,
	})

}

type FileInfo struct {
	Size             int64
	LastModification time.Time
}

func (c *Controller) UploadAjaxTriageFile(g *gin.Context) {
	// get requested path
	path := g.PostForm("path")

	// list folders in triage path
	entries, err := os.ReadDir(
		filepath.Join(c.config.Media.Path, TRIAGE_FILEPATH, path),
	)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	items := make(map[string]FileInfo)

	for _, entry := range entries {
		// check type
		entryPath := filepath.Join(
			c.config.Media.Path,
			TRIAGE_FILEPATH,
			path,
			entry.Name(),
		)
		stat, err := os.Stat(entryPath)
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		if stat.IsDir() {
			continue
		}

		items[entry.Name()] = FileInfo{
			Size:             stat.Size(),
			LastModification: stat.ModTime(),
		}
	}

	g.JSON(http.StatusOK, gin.H{
		"files": items,
	})

}

func (c *Controller) UploadAjaxUploadFile(g *gin.Context) {
	// get file
	file, err := g.FormFile("file")
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// get path
	_path := g.PostForm("path")
	path := filepath.Join(c.config.Media.Path, TRIAGE_FILEPATH, _path, file.Filename)

	// save file
	err = g.SaveUploadedFile(file, path)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{})
}
