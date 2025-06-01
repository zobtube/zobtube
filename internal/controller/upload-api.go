package controller

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

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

func (c *Controller) UploadAjaxDeleteFile(g *gin.Context) {
	// get file from request
	type fileDeleteForm struct {
		File string
	}

	form := fileDeleteForm{}
	err := g.ShouldBind(&form)
	if err != nil {
		g.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ensure not empty
	file := form.File
	if file == "" {
		g.JSON(400, gin.H{
			"error": "file name cannot be empty",
		})
		return
	}

	// assemble with triage path
	file = filepath.Join(c.config.Media.Path, TRIAGE_FILEPATH, file)

	// remove file
	err = os.Remove(file)
	if err != nil {
		g.JSON(422, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) UploadAjaxFolderCreate(g *gin.Context) {
	// get new folder name
	name := g.PostForm("name")

	// construct absolute path
	path := filepath.Join(c.config.Media.Path, TRIAGE_FILEPATH, name)

	// check if folder already exists
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		g.JSON(409, gin.H{
			"error": "Folder already exists",
		})
		return
	}

	// do not exists, create it
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{})
}
