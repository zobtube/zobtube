package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
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

func (c *Controller) UploadAjaxMassDelete(g *gin.Context) {
	// get file list from request
	type fileDeleteForm struct {
		Files []string `json:"files" binding:"required"`
	}

	form := fileDeleteForm{}
	err := g.ShouldBindJSON(&form)
	if err != nil {
		g.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ensure not empty
	files := form.Files
	if len(files) == 0 {
		g.JSON(400, gin.H{
			"error": "mass deletion requested without any files",
		})
		return
	}
	for _, file := range files {
		c.logger.Debug().Str("file", file).Send()
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
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) UploadAjaxMassImport(g *gin.Context) {
	// get file list from request
	type fileImportForm struct {
		Files      []string `json:"files" binding:"required"`
		Actors     []string `json:"actors"`
		Categories []string `json:"categories"`
		TypeEnum   string   `json:"type" binding:"required"`
		Channel    string   `json:"channel"`
	}

	form := fileImportForm{}
	err := g.ShouldBindJSON(&form)
	if err != nil {
		g.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ensure type is valid
	if form.TypeEnum != "c" && form.TypeEnum != "v" && form.TypeEnum != "m" {
		g.JSON(400, gin.H{
			"error": "type of video is invalid",
		})
	}

	// ensure not empty
	files := form.Files
	if len(files) == 0 {
		g.JSON(400, gin.H{
			"error": "mass import requested without any files",
		})
		return
	}

	// pre-check: ensure actors exists
	var actors []*model.Actor
	for _, actorID := range form.Actors {
		actor := &model.Actor{
			ID: actorID,
		}
		result := c.datastore.First(actor)

		// check result
		if result.RowsAffected < 1 {
			g.JSON(400, gin.H{
				"error": fmt.Sprintf("actor id %s does not exist", actorID),
			})
			return
		}

		actors = append(actors, actor)
	}

	// pre-check: ensure categories exists
	var categories []*model.CategorySub
	for _, subCategoryID := range form.Categories {
		subCategory := &model.CategorySub{
			ID: subCategoryID,
		}
		result := c.datastore.First(subCategory)

		// check result
		if result.RowsAffected < 1 {
			g.JSON(400, gin.H{
				"error": fmt.Sprintf("category id %s does not exist", subCategoryID),
			})
			return
		}

		categories = append(categories, subCategory)
	}

	// pre-check: ensure channel exists
	var channel *model.Channel
	channel = nil
	if form.Channel != "" {
		channel = &model.Channel{
			ID: form.Channel,
		}
		result := c.datastore.First(channel)

		// check result
		if result.RowsAffected < 1 {
			g.JSON(400, gin.H{
				"error": fmt.Sprintf("channel id %s does not exist", form.Channel),
			})
			return
		}
	}

	// pre-check: ensure files are not empty
	for _, file := range files {
		c.logger.Debug().Str("file", file).Send()
		if file == "" {
			g.JSON(400, gin.H{
				"error": "file name cannot be empty",
			})
			return
		}
	}

	// prepare transaction for the whole import
	tx := c.datastore.Begin()
	var videos []*model.Video

	// now perform file import
	for _, file := range files {
		video := &model.Video{
			Name:      file,
			Filename:  file,
			Type:      form.TypeEnum,
			Imported:  false,
			Thumbnail: false,
		}

		if channel != nil {
			video.Channel = channel
		}

		// save object in db
		err = tx.Debug().Create(&video).Error
		if err != nil {
			tx.Rollback()
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = tx.Model(video).Debug().Association("Actors").Append(actors)
		if err != nil {
			tx.Rollback()
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = tx.Model(video).Debug().Association("Categories").Append(categories)
		if err != nil {
			tx.Rollback()
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		videos = append(videos, video)
	}

	// validate transaction
	tx.Commit()

	// now create task for the import
	for _, video := range videos {
		err = c.runner.NewTask("video/create", map[string]string{
			"videoID":         video.ID,
			"thumbnailTiming": "0",
		})
		if err != nil {
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}
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
