package controller

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) CategorySubThumb(g *gin.Context) {
	// get id from path
	id := g.Param("id")

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

	// check if thumbnail exists
	if !category.Thumbnail {
		g.Redirect(http.StatusFound, CATEGORY_PROFILE_PICTURE_MISSING)
		return
	}

	// construct file path
	filename := fmt.Sprintf("%s.jpg", id)
	targetPath := filepath.Join(c.config.Media.Path, CATEGORY_FILEPATH, filename)

	// give file path
	g.File(targetPath)
}
