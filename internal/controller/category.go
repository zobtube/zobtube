package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

// CategoryList godoc
//
//	@Summary	List all categories
//	@Tags		category
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Router		/category [get]
func (c *Controller) CategoryList(g *gin.Context) {
	var categories []model.Category
	c.datastore.Preload("Sub").Find(&categories)
	g.JSON(http.StatusOK, gin.H{
		"items": categories,
		"total": len(categories),
	})
}

// CategorySubGet godoc
//
//	@Summary	Get category sub with videos and actors
//	@Tags		category
//	@Produce	json
//	@Param		id	path	string	true	"Category sub ID"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/category/{id} [get]
func (c *Controller) CategorySubGet(g *gin.Context) {
	id := g.Param("id")
	sub := &model.CategorySub{ID: id}
	result := c.datastore.Preload("Videos").Preload("Actors.Videos").First(sub)
	if result.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	g.JSON(http.StatusOK, sub)
}

// CategoryAdd godoc
//
//	@Summary	Create a new category
//	@Tags		category
//	@Accept		x-www-form-urlencoded
//	@Param		Name	formData	string	true	"Category name"
//	@Param		Type	formData	string	false	"Category type"
//	@Param		Scope	formData	string	false	"Category scope"
//	@Success	200
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/category [post]
func (c *Controller) CategoryAdd(g *gin.Context) {
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

// CategoryDelete godoc
//
//	@Summary	Delete a category
//	@Tags		category
//	@Param		id	path	string	true	"Category ID"
//	@Success	200
//	@Failure	400	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/category/{id} [delete]
func (c *Controller) CategoryDelete(g *gin.Context) {
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
