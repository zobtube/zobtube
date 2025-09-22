package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) CategoryList(g *gin.Context) {
	var categories []model.Category
	c.datastore.Preload("Sub").Find(&categories)
	c.HTML(g, http.StatusOK, "category/list.html", gin.H{
		"Categories": categories,
	})
}

func (c *Controller) CategorySubView(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	sub := &model.CategorySub{
		ID: id,
	}
	result := c.datastore.Preload("Videos").Preload("Actors.Videos").First(sub)

	// check result
	if result.RowsAffected < 1 {
		//TODO: return to homepage
		g.JSON(404, gin.H{})
		return
	}

	c.HTML(g, http.StatusOK, "category/view.html", gin.H{
		"Sub": sub,
	})
}
