package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) CategoryAjaxAdd(g *gin.Context) {
	var err error

	// check method
	if g.Request.Method != "POST" {
		g.JSON(405, gin.H{})
		return
	}

	// check form
	type CategoryForm struct {
		Name  string
		Type  string
		Scope string
	}
	var form CategoryForm
	err = g.ShouldBind(&form)
	if err != nil {
		g.JSON(400, gin.H{
			"error": err,
		})
		return
	}

	// check emptiness
	if form.Name == "" {
		g.JSON(400, gin.H{
			"error": "category name cannot be empty",
		})
		return
	}

	// create object
	category := &model.Category{
		Name: form.Name,
	}
	err = c.datastore.Create(&category).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) CategoryAjaxDelete(g *gin.Context) {
	// get category id from path
	id := g.Param("id")

	category := &model.Category{
		ID: id,
	}
	result := c.datastore.Preload("Sub").First(category)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	if len(category.Sub) > 0 {
		g.JSON(400, gin.H{
			"error": "category cannot be deleted with values presents",
		})
		return
	}

	err := c.datastore.Delete(&category).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error":       err.Error(),
			"human_error": "unable to delete category",
		})
		return
	}

	g.JSON(200, gin.H{})
}
