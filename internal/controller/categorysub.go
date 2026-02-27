package controller

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

// CategorySubAdd godoc
//
//	@Summary	Create a new category sub
//	@Tags		category
//	@Accept		x-www-form-urlencoded
//	@Param		Name	formData	string	true	"Sub-category name"
//	@Param		Parent	formData	string	true	"Parent category ID"
//	@Success	200
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/category-sub [post]
func (c *Controller) CategorySubAdd(g *gin.Context) {
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

// CategorySubThumbSet godoc
//
//	@Summary	Set category sub thumbnail
//	@Tags		category
//	@Accept		multipart/form-data
//	@Param		id	path	string	true	"Category sub ID"
//	@Param		pp	formData	file	true	"Thumbnail image"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/category-sub/{id}/thumb [post]
func (c *Controller) CategorySubThumbSet(g *gin.Context) {
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
		g.JSON(500, gin.H{"error": err.Error(), "human_error": "unable to retrieve thumbnail from form"})
		return
	}
	store, err := c.storageResolver.Storage(c.config.DefaultLibraryID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	filename := fmt.Sprintf("%s.jpg", category.ID)
	path := filepath.Join("categories", filename)
	src, err := file.Open()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()
	dst, err := store.Create(path)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
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

// CategorySubThumbRemove godoc
//
//	@Summary	Remove category sub thumbnail
//	@Tags		category
//	@Param		id	path	string	true	"Category sub ID"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/category-sub/{id}/thumb [delete]
func (c *Controller) CategorySubThumbRemove(g *gin.Context) {
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
	path := filepath.Join("categories", filename)
	store, err := c.storageResolver.Storage(c.config.DefaultLibraryID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	_ = store.Delete(path)

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

// CategorySubRename godoc
//
//	@Summary	Rename category sub
//	@Tags		category
//	@Accept		x-www-form-urlencoded
//	@Param		id	path	string	true	"Category sub ID"
//	@Param		title	formData	string	true	"New name"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/category-sub/{id}/rename [post]
func (c *Controller) CategorySubRename(g *gin.Context) {
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

// CategorySubThumb godoc
//
//	@Summary	Get category sub thumbnail image
//	@Tags		category
//	@Param		id	path	string	true	"Category sub ID"
//	@Success	200	file	bytes
//	@Failure	404
//	@Router		/category-sub/{id}/thumb [get]
func (c *Controller) CategorySubThumb(g *gin.Context) {
	id := g.Param("id")
	category := &model.CategorySub{ID: id}
	result := c.datastore.First(category)
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}
	if !category.Thumbnail {
		g.Redirect(http.StatusFound, CATEGORY_PROFILE_PICTURE_MISSING)
		return
	}
	store, err := c.storageResolver.Storage(c.config.DefaultLibraryID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	filename := fmt.Sprintf("%s.jpg", id)
	path := filepath.Join("categories", filename)
	c.serveFromStorage(g, store, path)
}
