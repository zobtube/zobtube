package controller

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) CategorySubAjaxAdd(g *gin.Context) {
	var err error
	if g.Request.Method == "POST" {
		type CategorySubForm struct {
			Name   string
			Parent string
		}
		var form CategorySubForm
		err = g.ShouldBind(&form)
		if err == nil {
			category := &model.CategorySub{
				Name:     form.Name,
				Category: form.Parent,
			}
			err = c.datastore.Create(&category).Error
			if err == nil {
				g.JSON(200, gin.H{})
				return
			}
		}
	}

	g.JSON(500, gin.H{
		"error": err,
	})
}

func (c *Controller) CategorySubAjaxThumbSet(g *gin.Context) {
	// get item from ID
	category := &model.CategorySub{
		ID: g.Param("id"),
	}
	result := c.datastore.First(category)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	file, err := g.FormFile("pp")
	if err != nil {
		g.JSON(500, gin.H{
			"error":       err.Error(),
			"human_error": "unable to retrieve thumbnail from form",
		})
		return
	}

	// construct file path
	filename := fmt.Sprintf("%s.jpg", category.ID)
	targetPath := filepath.Join(c.config.Media.Path, CATEGORY_FILEPATH, filename)

	// save thumb on disk
	err = g.SaveUploadedFile(file, targetPath)
	if err != nil {
		g.JSON(500, gin.H{
			"error":       err.Error(),
			"human_error": "unable to save uploaded thumbnail",
		})
		return
	}

	// check if thumbnail exists
	if !category.Thumbnail {
		category.Thumbnail = true
		err = c.datastore.Save(category).Error
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// all good
	g.JSON(200, gin.H{})
}

func (c *Controller) CategorySubAjaxThumbRemove(g *gin.Context) {
	// get item from ID
	category := &model.CategorySub{
		ID: g.Param("id"),
	}
	result := c.datastore.First(category)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// construct file path
	filename := fmt.Sprintf("%s.jpg", category.ID)
	targetPath := filepath.Join(c.config.Media.Path, CATEGORY_FILEPATH, filename)

	_, err := os.Stat(targetPath)
	if err != nil && !os.IsNotExist(err) {
		g.JSON(500, gin.H{
			"error":       err,
			"human_error": "unable to check video presence on disk",
		})
		return
	}
	if !os.IsNotExist(err) {
		// exist, deleting it
		err = os.Remove(targetPath)
		if err != nil {
			g.JSON(500, gin.H{
				"error":       err,
				"human_error": "unable to delete video on disk",
			})
		}
	}

	// check if thumbnail exists
	if category.Thumbnail {
		category.Thumbnail = false
		err = c.datastore.Save(category).Error
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// all good
	g.JSON(200, gin.H{})
}

func (c *Controller) CategorySubAjaxRename(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get new name
	var form struct {
		Title string `form:"title"`
	}
	err := g.ShouldBind(&form)
	if err != nil {
		// method not allowed
		g.JSON(406, gin.H{})
		return
	}

	// get item from ID
	category := &model.CategorySub{
		ID: id,
	}
	result := c.datastore.First(category)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	category.Name = form.Title

	err = c.datastore.Save(category).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// all good
	g.JSON(200, gin.H{})
}
